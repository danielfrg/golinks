<script lang="ts">
 import { Button, Table, TableHeaderCell } from "@kenn-io/kit-ui";
 import type { Link } from "$lib/links.js";
 let {links,publicUrl,readOnly,disabled=false,onEdit}:{links:Link[];publicUrl:string;readOnly:boolean;disabled?:boolean;onEdit:(link:Link)=>void}=$props();
</script>

<Table ariaLabel="Links" unsorted>
 {#snippet header()}
  <TableHeaderCell label="Short link"/>
  <TableHeaderCell label="Destination"/>
  <TableHeaderCell label="Visits" numeric/>
  {#if !readOnly}<TableHeaderCell><span class="sr-only">Actions</span></TableHeaderCell>{/if}
 {/snippet}
 {#each links as link (link.short)}
  <tr>
   <td><a class="short-link" href={`${publicUrl.replace(/\/$/,"")}/${encodeURIComponent(link.short)}`} target="_blank" rel="noreferrer"><span class="go-prefix">go/</span>{link.short}</a></td>
   <td><span class="destination" title={link.long}>{link.long}</span></td>
   <td><span class="visits" aria-label={`${link.clicks} visits`}>{link.clicks.toLocaleString()}</span></td>
   {#if !readOnly}<td><span id={`edit-${link.short}`}><Button class="row-action" disabled={disabled} onclick={()=>onEdit(link)} ariaLabel={`Edit ${link.short}`}>Edit</Button></span></td>{/if}
  </tr>
 {/each}
</Table>
