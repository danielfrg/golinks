<script lang="ts">
 import { Button } from "$lib/components/ui/button/index.js";
 import type { Link } from "$lib/links.js";
 let {links,publicUrl,readOnly,disabled=false,onEdit}:{links:Link[];publicUrl:string;readOnly:boolean;disabled?:boolean;onEdit:(link:Link)=>void}=$props();
</script>

<div class="link-list">
 <div class="list-labels" aria-hidden="true"><span>Short link</span><span>Destination</span><span class="visits-label">Visits</span>{#if !readOnly}<span></span>{/if}</div>
 <ul>
  {#each links as link (link.short)}
   <li class:read-only={readOnly}>
    <a class="short-link" href={`${publicUrl.replace(/\/$/,"")}/${encodeURIComponent(link.short)}`} target="_blank" rel="noreferrer"><span class="go-prefix">go/</span>{link.short}<svg aria-hidden="true" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M7 17 17 7M7 7h10v10" stroke-linecap="round" stroke-linejoin="round"/></svg></a>
    <span class="destination" title={link.long}>{link.long}</span>
    <span class="visits" aria-label={`${link.clicks} visits`}>{link.clicks.toLocaleString()}</span>
    {#if !readOnly}<Button id={`edit-${link.short}`} variant="ghost" class="row-action" disabled={disabled} onclick={()=>onEdit(link)} aria-label={`Edit ${link.short}`}>Edit</Button>{/if}
   </li>
  {/each}
 </ul>
</div>
