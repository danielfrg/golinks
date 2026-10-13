<script lang="ts">
 import {onMount} from "svelte";
 import {Button} from "$lib/components/ui/button/index.js";
 import {Input} from "$lib/components/ui/input/index.js";
 import {importLinks,type ImportResult} from "$lib/links.js";
 let {onImport,onCancel}:{onImport:(result:ImportResult)=>void;onCancel:()=>void}=$props();
 let linksFile=$state<File|null>(null);let statsFile=$state<File|null>(null);
 let fileInput:HTMLInputElement|null=$state(null);let busy=$state(false);let error=$state("");
 onMount(()=>fileInput?.focus());
 function selected(event:Event,kind:"links"|"stats"){
  const file=(event.currentTarget as HTMLInputElement).files?.[0]??null;
  if(kind==="links")linksFile=file;else statsFile=file;error="";
 }
 async function submit(event:SubmitEvent){
  event.preventDefault();if(busy||!linksFile)return;error="";
  if(linksFile.size+(statsFile?.size??0)>8*1024*1024){error="Choose files under 8 MiB combined, or use the command-line import for a larger backup.";return;}
  busy=true;
  try{const [links,stats]=await Promise.all([linksFile.text(),statsFile?.text()]);onImport(await importLinks(links,stats));}
  catch(e){error=e instanceof TypeError?"Couldn't reach Golinks. Try again.":e instanceof Error?e.message:"Couldn't import the backup. Try again.";}
  finally{busy=false;}
 }
</script>

<section class="editor" aria-labelledby="import-title">
 <div class="editor-heading"><h2 id="import-title">Import links</h2><Button variant="ghost" class="close-editor" onclick={onCancel} disabled={busy} aria-label="Close import"><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="m6 6 12 12M6 18 18 6" stroke-linecap="round"/></svg></Button></div>
 <p class="import-hint" id="import-hint">Choose a file saved with Export links. Existing links stay as they are.</p>
 <form onsubmit={submit}>
  <div class="editor-fields">
   <div class="field"><label for="import-links">Links file</label><Input id="import-links" bind:ref={fileInput} type="file" class="backup-file" accept=".jsonl,.ndjson,.json,application/x-ndjson,application/json" onchange={(event)=>selected(event,"links")} disabled={busy} required aria-describedby="import-hint"/></div>
   <div class="field"><label for="import-stats">Click history <span class="optional">(optional)</span></label><Input id="import-stats" type="file" class="backup-file" accept=".csv,text/csv" onchange={(event)=>selected(event,"stats")} disabled={busy}/><p class="field-hint">Without a CSV, total visits are kept. <a class="text-link" href="/.export-stats">Export click history</a></p></div>
  </div>
  {#if error}<p class="form-error" role="alert">{error}</p>{/if}
  <div class="editor-actions"><Button type="submit" class="control" disabled={busy||!linksFile}>{busy?"Importing…":"Import links"}</Button><Button variant="outline" class="control" onclick={onCancel} disabled={busy}>Cancel</Button></div>
 </form>
</section>
