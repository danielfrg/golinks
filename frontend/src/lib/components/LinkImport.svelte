<script lang="ts">
  import { onMount } from "svelte";
  import { XIcon } from "@lucide/svelte";
  import { Button } from "@ui-sdk/svelte/button";
  import { Field, FieldLabel } from "@ui-sdk/svelte/field";
  import { FileInput, FileInputTrigger } from "@ui-sdk/svelte/file-input";
  import { importLinks, type ImportResult } from "$lib/links.js";
  let { onImport, onCancel }: { onImport: (result: ImportResult) => void; onCancel: () => void } =
    $props();
  let linksFile = $state<File | null>(null);
  let statsFile = $state<File | null>(null);
  let busy = $state(false);
  let error = $state("");
  onMount(() => document.querySelector<HTMLElement>("[data-slot=file-input-trigger]")?.focus());
  function selected(files: File[], kind: "links" | "stats") {
    const file = files[0] ?? null;
    if (kind === "links") linksFile = file;
    else statsFile = file;
    error = "";
  }
  async function submit(event: SubmitEvent) {
    event.preventDefault();
    if (busy || !linksFile) return;
    error = "";
    if (linksFile.size + (statsFile?.size ?? 0) > 8 * 1024 * 1024) {
      error =
        "Choose files under 8 MiB combined, or use the command-line import for a larger backup.";
      return;
    }
    busy = true;
    try {
      const [links, stats] = await Promise.all([linksFile.text(), statsFile?.text()]);
      onImport(await importLinks(links, stats));
    } catch (e) {
      error =
        e instanceof TypeError
          ? "Couldn't reach Golinks. Try again."
          : e instanceof Error
            ? e.message
            : "Couldn't import the backup. Try again.";
    } finally {
      busy = false;
    }
  }
</script>

<section class="editor" aria-labelledby="import-title">
  <div class="editor-heading">
    <h2 id="import-title">Import links</h2>
    <Button
      class="close-editor"
      data-variant="ghost"
      onclick={onCancel}
      disabled={busy}
      aria-label="Close import"><XIcon aria-hidden="true" /></Button
    >
  </div>
  <p class="import-hint" id="import-hint">
    Choose a file saved with Export links. Existing links stay as they are.
  </p>
  <form onsubmit={submit}>
    <div class="editor-fields">
      <Field name="links">
        <FieldLabel>Links file</FieldLabel>
        <FileInput
          name="links"
          required
          disabled={busy}
          accept=".jsonl,.ndjson,.json,application/x-ndjson,application/json"
          onValueChange={(files) => selected(files, "links")}
        >
          <div class="file-picker">
            <FileInputTrigger aria-describedby="import-hint">Choose file</FileInputTrigger>
            <span class="file-name">{linksFile?.name ?? "No file chosen"}</span>
          </div>
        </FileInput>
      </Field>
      <Field name="stats">
        <FieldLabel>Click history <span class="optional">(optional)</span></FieldLabel>
        <FileInput
          name="stats"
          disabled={busy}
          accept=".csv,text/csv"
          onValueChange={(files) => selected(files, "stats")}
        >
          <div class="file-picker">
            <FileInputTrigger>Choose file</FileInputTrigger>
            <span class="file-name">{statsFile?.name ?? "No file chosen"}</span>
          </div>
        </FileInput>
      </Field>
    </div>
    {#if error}<p class="form-error" role="alert">{error}</p>{/if}
    <div class="editor-actions">
      <Button type="submit" data-variant="solid" disabled={busy || !linksFile}
        >{busy ? "Importing…" : "Import links"}</Button
      ><Button onclick={onCancel} disabled={busy}>Cancel</Button>
    </div>
  </form>
</section>
