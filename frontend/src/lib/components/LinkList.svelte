<script lang="ts">
  import { Button } from "@ui-sdk/svelte/button";
  import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
  } from "../../design-system/custom/table";
  import type { Link } from "$lib/links.js";
  let {
    links,
    publicUrl,
    readOnly,
    disabled = false,
    onEdit,
  }: {
    links: Link[];
    publicUrl: string;
    readOnly: boolean;
    disabled?: boolean;
    onEdit: (link: Link) => void;
  } = $props();
</script>

<Table aria-label="Links">
  <TableHeader>
    <TableRow>
      <TableHead>Short link</TableHead>
      <TableHead>Destination</TableHead>
      <TableHead align="end">Visits</TableHead>
      {#if !readOnly}<TableHead align="end"><span class="sr-only">Actions</span></TableHead>{/if}
    </TableRow>
  </TableHeader>
  <TableBody>
    {#each links as link (link.short)}
      <TableRow>
        <TableCell
          ><a
            class="short-link"
            href={`${publicUrl.replace(/\/$/, "")}/${encodeURIComponent(link.short)}`}
            target="_blank"
            rel="noreferrer"><span class="go-prefix">go/</span>{link.short}</a
          ></TableCell
        >
        <TableCell><span class="destination" title={link.long}>{link.long}</span></TableCell>
        <TableCell align="end"
          ><span class="visits" aria-label={`${link.clicks} visits`}
            >{link.clicks.toLocaleString()}</span
          ></TableCell
        >
        {#if !readOnly}<TableCell align="end"
            ><span id={`edit-${link.short}`}
              ><Button {disabled} onclick={() => onEdit(link)} aria-label={`Edit ${link.short}`}
                >Edit</Button
              ></span
            ></TableCell
          >{/if}
      </TableRow>
    {/each}
  </TableBody>
</Table>
