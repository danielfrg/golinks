<script lang="ts">
 import { onMount, untrack } from "svelte";
 import { Button, IconButton } from "@kenn-io/kit-ui";
 import { TextInput } from "@kenn-io/kit-ui";
 import {createLink,updateLink,deleteLink,previewLink,type Link} from "$lib/links.js";
 let {link=null,suggested="",onSave,onDelete,onCancel}:{link?:Link|null;suggested?:string;onSave:(link:Link)=>void;onDelete:(short:string)=>void;onCancel:()=>void}=$props();
 let short=$state(untrack(()=>link?.short??suggested));
 let destination=$state(untrack(()=>link?.long??""));
 let busy=$state(false);let error=$state("");let deleting=$state(false);
 let testing=$state(false);let path=$state("");let query=$state("");let preview=$state("");let previewError=$state("");let previewing=$state(false);
 let nameInput:HTMLInputElement|undefined=$state();let destinationInput:HTMLTextAreaElement|null=$state(null);
 let previewRequest:AbortController|null=null;
 onMount(()=>{if(link)destinationInput?.focus();else nameInput?.focus();return()=>previewRequest?.abort();});
 function invalidatePreview(){previewRequest?.abort();preview="";previewError="";previewing=false;}
 async function save(event:SubmitEvent){event.preventDefault();if(busy)return;busy=true;error="";
  try{const result=link?await updateLink(link.short,destination.trim()):await createLink(short.trim(),destination.trim());onSave(result);}
  catch(e){error=e instanceof Error?e.message:"Couldn't save the link. Try again.";}finally{busy=false;}
 }
 async function remove(){if(!link||busy)return;busy=true;error="";try{await deleteLink(link.short);onDelete(link.short);}catch(e){error=e instanceof Error?e.message:"Couldn't delete the link. Try again.";}finally{busy=false;}}
 async function tryPath(){invalidatePreview();const request=new AbortController();previewRequest=request;previewing=true;
  try{const result=await previewLink(destination,path,query,request.signal);if(!request.signal.aborted)preview=result.url;}
  catch(e){if(!request.signal.aborted)previewError=e instanceof Error?e.message:"Couldn't preview this destination.";}
  finally{if(previewRequest===request)previewing=false;}
 }
</script>

<section class="editor" aria-labelledby="editor-title">
 <div class="editor-heading"><h2 id="editor-title">{link?"Edit link":"New link"}</h2><IconButton class="close-editor" onclick={onCancel} disabled={busy} ariaLabel="Close editor"><svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="m6 6 12 12M6 18 18 6" stroke-linecap="round"/></svg></IconButton></div>
 <form onsubmit={save}>
  <div class="editor-fields">
   <div class="field"><label for="short-name">Short name</label><TextInput block id="short-name" bind:inputEl={nameInput} bind:value={short} disabled={busy||!!link} required autocomplete="off" placeholder="meeting-notes" ariaDescribedby="name-hint"/><p class="field-hint" id="name-hint">Case and hyphens are ignored.</p></div>
   <div class="field"><label for="destination">Destination</label><textarea id="destination" bind:this={destinationInput} bind:value={destination} oninput={invalidatePreview} disabled={busy} required rows="2" placeholder="https://example.com/" spellcheck="false" aria-describedby={error?"save-error":undefined} aria-invalid={error?true:undefined}></textarea></div>
  </div>
  {#if error}<p class="form-error" id="save-error" role="alert">{error}</p>{/if}
  <div class="editor-actions"><Button type="submit" class="control" disabled={busy}>{busy?"Saving…":"Save link"}</Button><Button class="control" onclick={onCancel} disabled={busy}>Cancel</Button><Button class="test-toggle" onclick={()=>{testing=!testing;}} disabled={busy} ariaExpanded={testing}>Try a path<svg class:expanded={testing} aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="m8 10 4 4 4-4" stroke-linecap="round" stroke-linejoin="round"/></svg></Button></div>
  {#if testing}<div class="preview-area" id="path-preview">
   <div class="preview-fields"><div class="field"><label for="preview-path">Extra path</label><TextInput block id="preview-path" bind:value={path} oninput={invalidatePreview} placeholder="docs/getting-started" disabled={busy}/></div><div class="field"><label for="preview-query">Query string <span class="optional">(optional)</span></label><TextInput block id="preview-query" bind:value={query} oninput={invalidatePreview} placeholder="page=2" disabled={busy}/></div><Button class="control preview-button" onclick={tryPath} disabled={busy||previewing||!destination}>{previewing?"Testing…":"Preview"}</Button></div>
   {#if previewError}<p class="form-error" role="alert">{previewError}</p>{/if}
   {#if preview}<div class="preview-result" role="status"><span>Redirects to</span><code>{preview}</code></div>{/if}
  </div>{/if}
 </form>
 {#if link}<div class="delete-area">{#if deleting}<p>Delete <strong>{link.short}</strong>?</p><Button tone="danger" class="control" onclick={remove} disabled={busy}>Delete link</Button><Button class="control" onclick={()=>deleting=false} disabled={busy}>Keep link</Button>{:else}<Button tone="danger" class="delete-toggle" onclick={()=>deleting=true} disabled={busy}>Delete link</Button>{/if}</div>{/if}
</section>
