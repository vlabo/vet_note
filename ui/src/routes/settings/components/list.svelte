<script lang="ts">
  import { onMount } from "svelte";
  import { writable } from "svelte/store";
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
  import { deleteSetting } from "$lib/DataService";
  import { dndzone, SOURCES, TRIGGERS } from "svelte-dnd-action";

  export let items = writable<Array<ViewSetting>>([]);
  export let title: string;
  export let placeholder: string;
  export let settingType: "PatientType" | "ProcedureType" | "PatientFolder" =
    "PatientType";
  export let isAddingItem = false;

  let deletePopupIndex = -1;
  let newItem = "";
  let dragDisabled = true;

  function addItem() {
    isAddingItem = true;
  }

  function acceptItem() {
    console.log(settingType);
    items.update(($items: ViewSetting[]): ViewSetting[] => {
      $items.push({
        id: undefined,
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

  function handleDnd(e: any) {
    const {
      items: newItems,
      info: { source, trigger },
    } = e.detail;
    items.update(() => {
      newItems.forEach((s: ViewSetting, index: number) => {
        s.index = index;
      });
      return newItems;
    });
    if (source === SOURCES.KEYBOARD && trigger === TRIGGERS.DRAG_STOPPED) {
      dragDisabled = true;
    }
  }

  function handleDndFinalize(e: any) {
    const {
      items: _,
      info: { source },
    } = e.detail;
    handleDnd(e);
    onUpdate($items);

    if (source === SOURCES.POINTER) {
      dragDisabled = true;
    }
  }

  export var onUpdate = (_: ViewSetting[]) => {};

  function deleteItem() {
    items.update(($items: ViewSetting[]): ViewSetting[] => {
      const deleted = $items.splice(deletePopupIndex, 1);
      if (deleted.length > 0) {
        deleteSetting(deleted![0].id!);
      }
      deletePopupIndex = -1;
      $items.sort((a, b) => a.index! - b.index!);
      return $items;
    });
  }

  function startDrag(e: any) {
    // preventing default to prevent lag on touch devices (because of the browser checking for screen scrolling)
    e.preventDefault();
    dragDisabled = false;
  }
  function handleKeyDown(e: any) {
    if ((e.key === "Enter" || e.key === " ") && dragDisabled)
      dragDisabled = false;
  }
</script>

<div class="mt-4 mb-6">
  <h2 class="text-gray-700 font-medium border-b border-gray-400">{title}</h2>
  <section
    use:dndzone={{ items: $items, dragDisabled, type: settingType }}
    on:consider={handleDnd}
    on:finalize={handleDndFinalize}
  >
    {#each $items as item, index (item.id)}
      <div class="flex items-center border-b border-gray-400 bg-white">
        <div
          tabindex="0"
          role="button"
          aria-label="drag-handle"
          style={dragDisabled ? "cursor: grab" : "cursor: grabbing"}
          on:mousedown={startDrag}
          on:touchstart={startDrag}
          on:keydown={handleKeyDown}
        >
          <Icon icon={faBars} class="py-3 pl-2 pr-10 drag-handle" />
        </div>
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
