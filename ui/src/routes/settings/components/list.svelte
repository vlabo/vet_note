<script lang="ts">
  import { onMount } from "svelte";
  import { writable } from "svelte/store";
  import Sortable from "sortablejs";
  import Icon from "$lib/icon.svelte";
  import {
    faBars,
    faTrash,
    faPlus,
    faXmark,
    faCheck,
  } from "@fortawesome/free-solid-svg-icons";
  import type { ViewSetting } from "$lib/types";
  import DeletePopup from "$lib/DeletePopup.svelte";

  export let items = writable<Array<ViewSetting>>([]);
  export let title: string;
  export let placeholder: string;
  export const settingType: "PatientType" | "ProcedureType" = "PatientType";
  export let isAddingItem = false;

  let deletePopupIndex = -1;
  let newItem = "";
  let listElement: HTMLElement;

  onMount(() => {
    items.update(($items: ViewSetting[]): ViewSetting[] => {
      $items.sort((a, b) => a.index - b.index);
      return $items;
    });

    Sortable.create(listElement, {
      animation: 100,
      handle: ".drag-handle",
      onStart(evt) {
        evt.item.classList.add("bg-blue-300");
        evt.item.classList.remove("bg-white");
      },
      onEnd(evt) {
        evt.item.classList.remove("bg-blue-300");
        evt.item.classList.add("bg-white");
        items.update(($items: ViewSetting[]): ViewSetting[] => {
          const clonedItems = $items.slice();
          const movedItem = clonedItems.splice(evt.oldIndex!, 1)[0];
          clonedItems.splice(evt.newIndex!, 0, movedItem);
          clonedItems.forEach((item, index) => {
            item.index = index;
          });
          onUpdate(clonedItems);
          return $items;
        });
      },
      onClone(evt) {
        evt.item.classList.add("bg-blue-300");
        evt.item.classList.remove("bg-white");
      },
    });
  });

  function addItem() {
    isAddingItem = true;
  }

  function acceptItem() {
    items.update(($items: ViewSetting[]): ViewSetting[] => {
      $items.push({
        value: newItem,
        type: settingType,
        index: $items.length,
      });
      isAddingItem = false;
      newItem = "";
      onUpdate($items);
      return $items;
    });
  }

  export var onUpdate = (_: ViewSetting[]) => {};

  function deleteItem() {
    items.update(($items: ViewSetting[]): ViewSetting[] => {
      $items.splice(deletePopupIndex, 1);
      deletePopupIndex = -1;
      $items.sort((a, b) => a.index - b.index);
      return $items;
    });
  }
</script>

<div class="mt-4 mb-6">
  <h2 class="text-gray-700 font-medium border-b border-gray-400">{title}</h2>
  <section bind:this={listElement}>
    {#each $items as item, index}
      <div class="flex items-center border-b border-gray-400 bg-white">
        <Icon icon={faBars} class="py-3 pl-2 pr-10 drag-handle" />
        <span class="text-gray-800">{item.value}</span>
        <button
          on:click={() => (deletePopupIndex = index)}
          class="focus:outline-none p-2 ml-auto hover:bg-blue-200"
        >
          <Icon icon={faTrash} />
        </button>
      </div>
    {/each}
  </section>

  <div class="mt-4 mb-6">
    {#if isAddingItem}
      <div class="flex items-center space-x-6 p-2">
        <input
          bind:value={newItem}
          type="text"
          {placeholder}
          class="flex-1 border rounded px-2 py-2 focus:outline-none focus:ring"
        />
        <button
          class="bg-red-500 text-white px-4 py-1 rounded focus:outline-none"
          on:click={() => (isAddingItem = false)}
        >
          <Icon icon={faXmark} />
        </button>
        <button
          class="bg-green-500 text-white px-4 py-1 border rounded focus:outline-none"
          on:click={acceptItem}
        >
          <Icon icon={faCheck} />
        </button>
      </div>
    {:else}
      <button
        on:click={addItem}
        class="w-full border-2 border-dashed border-gray-300 p-2 rounded flex justify-center items-center focus:outline-none"
        ><Icon icon={faPlus} class="mx-2" /> Добави нов {title.toLowerCase()}</button
      >
    {/if}
  </div>
</div>

{#if deletePopupIndex != -1}
  <DeletePopup
    message={`Сигурни ли сте че искате да изтриете "${$items[deletePopupIndex].value}"?`}
    confirmAction={deleteItem}
    onComplete={() => {
      deletePopupIndex = -1;
    }}
  />
{/if}
