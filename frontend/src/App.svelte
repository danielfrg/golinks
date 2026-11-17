<script lang="ts">
  import { onMount, tick } from "svelte";
  import { CheckIcon, ChevronDownIcon, PlusIcon, SearchIcon } from "@lucide/svelte";
  import { Button } from "@ui-sdk/svelte/button";
  import { Input, InputGroup } from "@ui-sdk/svelte/input";
  import { Preferences } from "@ui-sdk/svelte/preferences";
  import {
    Select,
    SelectContent,
    SelectItem,
    SelectItemIndicator,
    SelectLabel,
    SelectTrigger,
    SelectValue,
  } from "@ui-sdk/svelte/select";
  import { ThemeToggle } from "./design-system/custom/theme-toggle";
  import { colorScheme } from "$lib/color-scheme";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import LinkList from "$lib/components/LinkList.svelte";
  import LinkEditor from "$lib/components/LinkEditor.svelte";
  import TemplateHelp from "$lib/components/TemplateHelp.svelte";
  import LinkImport from "$lib/components/LinkImport.svelte";
  import {
    fetchLinks,
    fetchSettings,
    filterLinks,
    sortLinks,
    readView,
    type Link,
    type Settings,
    type View,
    type ImportResult,
  } from "$lib/links.js";
  let links = $state<Link[]>([]);
  let settings = $state<Settings | null>(null);
  let loaded = $state(false);
  let loading = $state(false);
  let error = $state("");
  let status = $state("");
  const initialView = readView(new URL(window.location.href));
  let view = $state(initialView);
  let query = $state(initialView.query);
  let sort = $state<"name" | "clicks">("clicks");
  let request: AbortController | null = null;

  let importing = $state(false);
  let importStatus = $state("");
  const sortItems = [
    { value: "name", label: "Name" },
    { value: "clicks", label: "Most used" },
  ];
  const visible = $derived(sortLinks(filterLinks(links, query), sort));
  const normalize = (value: string) => value.toLowerCase().replaceAll("-", "");
  const editing = $derived(
    view.edit === null
      ? null
      : (links.find((link) => normalize(link.short) === normalize(view.edit!)) ?? null),
  );
  const editorOpen = $derived(
    settings !== null && !settings.readOnly && (view.create !== null || editing !== null),
  );

  function setView(patch: Partial<View>, replace = false) {
    if (patch.edit !== undefined || patch.create !== undefined) importing = false;
    view = { ...view, ...patch };
    query = view.query;
    const url = new URL("/", window.location.origin);
    if (view.query) url.searchParams.set("q", view.query);
    if (view.edit !== null) url.searchParams.set("edit", view.edit);
    if (view.create !== null) {
      url.searchParams.set("new", "1");
      if (view.create) url.searchParams.set("short", view.create);
    }
    if (view.help) url.searchParams.set("help", "1");
    window.history[replace ? "replaceState" : "pushState"]({}, "", url);
  }
  async function load() {
    if (loading) return;
    loading = true;
    error = "";
    const controller = new AbortController();
    request = controller;
    try {
      const [catalog, config] = await Promise.all([
        fetchLinks(controller.signal),
        fetchSettings(controller.signal),
      ]);
      links = catalog.links;
      settings = config;
      loaded = true;
    } catch (e) {
      if (!controller.signal.aborted)
        error = loaded
          ? "Couldn't refresh links. Showing the last loaded list."
          : e instanceof Error
            ? e.message
            : "Couldn't load links. Try again.";
    } finally {
      if (!controller.signal.aborted) loading = false;
    }
  }
  async function closeEditor(focusShort?: string) {
    setView({ edit: null, create: null });
    await tick();
    const target = focusShort ? document.getElementById(`edit-${focusShort}`) : null;
    (
      target?.querySelector<HTMLButtonElement>("button") ??
      document.querySelector<HTMLButtonElement>(".new-link")
    )?.focus();
  }
  function saved(link: Link) {
    links = [...links.filter((item) => normalize(item.short) !== normalize(link.short)), link];
    status = `${link.short} saved.`;
    void closeEditor(link.short);
  }
  function deleted(short: string) {
    links = links.filter((link) => link.short !== short);
    status = `${short} deleted.`;
    void closeEditor();
  }
  function openImport() {
    setView({ edit: null, create: null });
    importStatus = "";
    importing = true;
  }
  async function closeImport(reload = false) {
    importing = false;
    if (reload) await load();
    await tick();
    document.getElementById("import-action")?.focus();
  }
  function imported(result: ImportResult) {
    importStatus = `Imported ${result.imported} ${result.imported === 1 ? "link" : "links"}. Skipped ${result.skipped} existing ${result.skipped === 1 ? "link" : "links"}.`;
    if (result.skippedStats)
      importStatus += ` Skipped ${result.skippedStats} history rows for links outside the backup.`;
    void closeImport(true);
  }
  onMount(() => {
    void load();
    const onPop = () => {
      importing = false;
      view = readView(new URL(window.location.href));
      query = view.query;
    };
    window.addEventListener("popstate", onPop);
    return () => {
      request?.abort();
      window.removeEventListener("popstate", onPop);
    };
  });
</script>

<Preferences preferences={[colorScheme]}>
  <header class="app-header">
    <a
      class="app-title"
      href="/"
      onclick={(event) => {
        event.preventDefault();
        setView({ query: "", edit: null, create: null, help: false });
      }}><h1>Golinks</h1></a
    ><ThemeToggle {colorScheme} />
  </header>
  <main class="app-shell">
    {#if error}<div class="notice" role="alert">
        <p>{error}</p>
        <Button onclick={load} disabled={loading}>Try again</Button>
      </div>{/if}
    <p class="sr-only" role="status" aria-live="polite">{status}</p>
    {#if importStatus}<p class="import-status" role="status">{importStatus}</p>{/if}
    {#if importing && settings && !settings.readOnly}<LinkImport
        onImport={imported}
        onCancel={() => void closeImport()}
      />{/if}
    {#if editorOpen}{#key editing?.short ?? `new-${view.create}`}<LinkEditor
          link={editing}
          suggested={view.create ?? ""}
          onSave={saved}
          onDelete={deleted}
          onCancel={() => void closeEditor()}
        />{/key}
    {:else if loaded && view.edit !== null && !editing}<p class="notice" role="alert">
        That link no longer exists. <button class="text-link" onclick={() => void closeEditor()}
          >Back to links</button
        >
      </p>{/if}
    <section class="directory" aria-label="Links" aria-busy={loading}>
      <div class="directory-toolbar">
        <InputGroup class="search-field">
          <SearchIcon aria-hidden="true" />
          <Input
            type="search"
            aria-label="Search links"
            value={query}
            onValueChange={(value) => setView({ query: value }, true)}
            placeholder="Search links"
          />
        </InputGroup>
        {#if !settings?.readOnly}<Button
            class="new-link"
            data-variant="solid"
            onclick={() => {
              status = "";
              setView({ edit: null, create: "" });
            }}
            disabled={!settings || editorOpen || importing}
            ><PlusIcon aria-hidden="true" />New link</Button
          >{/if}
        {#if settings?.readOnly}<span>Read-only</span>{/if}
        <div class="sort-label">
          <Select
            items={sortItems}
            value={sort}
            onValueChange={(value) => (sort = value === "clicks" ? "clicks" : "name")}
          >
            <SelectLabel>Sort</SelectLabel>
            <SelectTrigger>
              <SelectValue />
              <ChevronDownIcon aria-hidden="true" />
            </SelectTrigger>
            <SelectContent>
              {#each sortItems as item (item.value)}
                <SelectItem value={item.value}>
                  {item.label}
                  <SelectItemIndicator><CheckIcon aria-hidden="true" /></SelectItemIndicator>
                </SelectItem>
              {/each}
            </SelectContent>
          </Select>
        </div>
      </div>
      {#if !loaded && loading}<div class="empty-state" role="status"><p>Loading links…</p></div>
      {:else if loaded && links.length === 0}<EmptyState
          title="No links yet"
          description={settings?.readOnly
            ? "No links have been added to this directory."
            : "Create a short name for a URL you use often."}
        />
      {:else if loaded && visible.length === 0}<EmptyState
          title="No matching links"
          description="Try another name or destination."
          ><Button onclick={() => setView({ query: "" }, true)}>Clear search</Button></EmptyState
        >
      {:else if loaded && settings}<LinkList
          links={visible}
          publicUrl={settings.publicUrl}
          readOnly={settings.readOnly}
          disabled={importing}
          onEdit={(link) => {
            status = "";
            setView({ edit: link.short, create: null });
          }}
        />{/if}
    </section>
    <footer class="app-footer">
      <Button
        class="footer-action"
        data-variant="link"
        onclick={() => setView({ help: !view.help })}
        aria-expanded={view.help}>Template help</Button
      ><a href="/_/export">Export links</a><a href="/_/export-stats">Export history</a
      >{#if settings && !settings.readOnly}<Button
          class="footer-action"
          data-variant="link"
          id="import-action"
          onclick={openImport}
          aria-expanded={importing}
          disabled={importing || editorOpen || loading}>Import links</Button
        >{/if}<a href="/_/api/docs">API docs</a>
    </footer>
    {#if view.help}<TemplateHelp />{/if}
  </main>
</Preferences>
