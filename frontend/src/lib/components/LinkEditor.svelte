<script lang="ts">
  import { onMount, untrack } from "svelte";
  import { ChevronDownIcon, XIcon } from "@lucide/svelte";
  import { Button } from "@ui-sdk/svelte/button";
  import { Field, FieldDescription, FieldLabel } from "@ui-sdk/svelte/field";
  import { Input } from "@ui-sdk/svelte/input";
  import { Textarea } from "@ui-sdk/svelte/textarea";
  import { createLink, updateLink, deleteLink, previewLink, type Link } from "$lib/links.js";
  let {
    link = null,
    suggested = "",
    onSave,
    onDelete,
    onCancel,
  }: {
    link?: Link | null;
    suggested?: string;
    onSave: (link: Link) => void;
    onDelete: (short: string) => void;
    onCancel: () => void;
  } = $props();
  let short = $state(untrack(() => link?.short ?? suggested));
  let destination = $state(untrack(() => link?.long ?? ""));
  let busy = $state(false);
  let error = $state("");
  let deleting = $state(false);
  let testing = $state(false);
  let path = $state("");
  let query = $state("");
  let preview = $state("");
  let previewError = $state("");
  let previewing = $state(false);
  let previewRequest: AbortController | null = null;
  onMount(() => {
    document.querySelector<HTMLElement>(link ? "[name=destination]" : "[name=short]")?.focus();
    return () => previewRequest?.abort();
  });
  function invalidatePreview() {
    previewRequest?.abort();
    preview = "";
    previewError = "";
    previewing = false;
  }
  async function save(event: SubmitEvent) {
    event.preventDefault();
    if (busy) return;
    busy = true;
    error = "";
    try {
      const result = link
        ? await updateLink(link.short, destination.trim())
        : await createLink(short.trim(), destination.trim());
      onSave(result);
    } catch (e) {
      error = e instanceof Error ? e.message : "Couldn't save the link. Try again.";
    } finally {
      busy = false;
    }
  }
  async function remove() {
    if (!link || busy) return;
    busy = true;
    error = "";
    try {
      await deleteLink(link.short);
      onDelete(link.short);
    } catch (e) {
      error = e instanceof Error ? e.message : "Couldn't delete the link. Try again.";
    } finally {
      busy = false;
    }
  }
  async function tryPath() {
    invalidatePreview();
    const request = new AbortController();
    previewRequest = request;
    previewing = true;
    try {
      const result = await previewLink(destination, path, query, request.signal);
      if (!request.signal.aborted) preview = result.url;
    } catch (e) {
      if (!request.signal.aborted)
        previewError = e instanceof Error ? e.message : "Couldn't preview this destination.";
    } finally {
      if (previewRequest === request) previewing = false;
    }
  }
</script>

<section class="editor" aria-labelledby="editor-title">
  <div class="editor-heading">
    <h2 id="editor-title">{link ? "Edit link" : "New link"}</h2>
    <Button
      class="close-editor"
      data-variant="ghost"
      onclick={onCancel}
      disabled={busy}
      aria-label="Close editor"><XIcon aria-hidden="true" /></Button
    >
  </div>
  <form onsubmit={save}>
    <div class="editor-fields">
      <Field name="short">
        <FieldLabel>Short name</FieldLabel>
        <Input
          bind:value={short}
          disabled={busy || !!link}
          required
          autocomplete="off"
          placeholder="meeting-notes"
        />
        <FieldDescription>Case and hyphens are ignored.</FieldDescription>
      </Field>
      <Field name="destination" invalid={error ? true : undefined}>
        <FieldLabel>Destination</FieldLabel>
        <Textarea
          bind:value={destination}
          oninput={invalidatePreview}
          disabled={busy}
          required
          rows={2}
          placeholder="https://example.com/"
          spellcheck="false"
          aria-describedby={error ? "save-error" : undefined}
        />
      </Field>
    </div>
    {#if error}<p class="form-error" id="save-error" role="alert">{error}</p>{/if}
    <div class="editor-actions">
      <Button type="submit" data-variant="solid" disabled={busy}
        >{busy ? "Saving…" : "Save link"}</Button
      ><Button onclick={onCancel} disabled={busy}>Cancel</Button><Button
        class="test-toggle"
        onclick={() => {
          testing = !testing;
        }}
        disabled={busy}
        aria-expanded={testing}
        aria-controls="path-preview"
        >Try a path<ChevronDownIcon class={testing ? "expanded" : ""} aria-hidden="true" /></Button
      >
    </div>
    {#if testing}<div class="preview-area" id="path-preview">
        <div class="preview-fields">
          <Field name="preview-path">
            <FieldLabel>Extra path</FieldLabel>
            <Input
              bind:value={path}
              oninput={invalidatePreview}
              placeholder="docs/getting-started"
              disabled={busy}
            />
          </Field>
          <Field name="preview-query">
            <FieldLabel>Query string <span class="optional">(optional)</span></FieldLabel>
            <Input
              bind:value={query}
              oninput={invalidatePreview}
              placeholder="page=2"
              disabled={busy}
            />
          </Field>
          <Button
            class="preview-button"
            onclick={tryPath}
            disabled={busy || previewing || !destination}
            >{previewing ? "Testing…" : "Preview"}</Button
          >
        </div>
        {#if previewError}<p class="form-error" role="alert">{previewError}</p>{/if}
        {#if preview}<div class="preview-result" role="status">
            <span>Redirects to</span><code>{preview}</code>
          </div>{/if}
      </div>{/if}
  </form>
  {#if link}<div class="delete-area">
      {#if deleting}<p>Delete <strong>{link.short}</strong>?</p>
        <Button data-variant="destructive" onclick={remove} disabled={busy}>Delete link</Button
        ><Button onclick={() => (deleting = false)} disabled={busy}>Keep link</Button>{:else}<Button
          data-variant="destructive"
          onclick={() => (deleting = true)}
          disabled={busy}>Delete link</Button
        >{/if}
    </div>{/if}
</section>
