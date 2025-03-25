<script lang="ts">
  import { type ViewProcedure } from "$lib/types";
  import { onMount } from "svelte";
  import { page } from "$app/state";
  import BackButton from "$lib/BackButton.svelte";
  import { faEdit, faTrash } from "@fortawesome/free-solid-svg-icons";
  import Icon from "$lib/icon.svelte";
  import DeletePopup from "$lib/DeletePopup.svelte";
  import {
    addProcedure,
    deleteProcedure,
    getProcedure,
    updateProcedure,
  } from "$lib/DataService";
  import { settingsProcedures } from "$lib/DataService";

  var procedure: ViewProcedure = {
    patientId: Number(page.params.id),
  };

  function formatDate(date: Date) {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0"); // getMonth() returns 0-based month
    const day = String(date.getDate()).padStart(2, "0");

    return `${year}-${month}-${day}`;
  }

  var date = formatDate(new Date());
  export var mode: "view" | "edit" | "new" = "view";
  var deletePopupOpen = false;

  onMount(() => {
    if (page.params.procedureId === "new") {
      mode = "new";
    } else {
      getProcedure(
        Number(page.params.id),
        Number(page.params.procedureId),
      ).then((p) => {
        if (p) {
          procedure = p;
          date = formatDate(new Date(p.date!));
        }
      });
    }
  });

  function onSave() {
    procedure.date = date;
    if (mode === "new") {
      addProcedure(Number(page.params.id), procedure);
    } else if (mode === "edit") {
      updateProcedure(Number(page.params.id), procedure);
    }
    history.back();
  }

  function onDelete() {
    deleteProcedure(Number(page.params.id), Number(page.params.procedureId));
    history.go(-2);
  }

  var types = settingsProcedures;
</script>

<div class="flex-col h-screen">
  <!-- Header -->
  <header
    class="bg-white shadow p-4 flex max-w-7xl mx-auto items-center justify-between"
  >
    <div>
      <BackButton />
    </div>
    <div>
      <h1 class="text-lg font-semibold">
        {#if mode === "view"}
          <span>Процедура</span>
        {:else if mode === "edit"}
          <span>Редактиране на Процедура</span>
        {:else if mode === "new"}
          <span>Нова Процедура</span>
        {/if}
      </h1>
    </div>
    <div class="flex space-x-2">
      {#if mode === "view"}
        <a
          href={`/patient/${page.params.id}/procedure/${page.params.procedureId}/edit`}
          class="text-black focus:outline-none mx-1 px-4 py-1 hover:bg-blue-100 hover:rounded-lg"
        >
          <Icon icon={faEdit} />
        </a>
      {:else if mode === "edit"}
        <button
          class="text-red-600 focus:outline-none mx-1 px-4 py-1 hover:bg-blue-100 hover:rounded-lg"
          on:click={() => {
            deletePopupOpen = true;
          }}
        >
          <Icon icon={faTrash} />
        </button>
      {/if}
    </div>
  </header>

  <!-- Form Content -->
  <main class="p-4 max-w-7xl mx-auto flex flex-col h-[calc(100vh-80px)]">
    <form class="flex flex-col flex-grow">
      <div class="flex flex-col space-y-2 flex-grow">
        <!-- Procedure type -->
        {#if mode === "view"}
          <span class="text-gray-800">{procedure.type}</span>
        {:else}
          <select
            bind:value={procedure.type}
            class="border border-gray-300 rounded p-2"
          >
            {#each $types as type}
              <option>{type.value}</option>
            {/each}
          </select>
        {/if}
        <!-- Datepicker -->
        <input
          type="date"
          disabled={mode == "view"}
          bind:value={date}
          class="border border-gray-300 rounded p-2"
          data-datepicker
        />

        <!-- Procedure details -->
        <textarea
          bind:value={procedure.details}
          disabled={mode == "view"}
          placeholder="Описание"
          class="w-full border border-gray-300 rounded p-2 resize-none flex-grow"
        ></textarea>
        <!-- Save button for edit/new mode only -->
        {#if mode !== "view"}
          <button
            type="submit"
            class="w-full bg-blue-600 text-white p-2 rounded hover:bg-blue-700 transition"
            on:click={onSave}
          >
            Запиши
          </button>
        {/if}
      </div>
    </form>
  </main>
</div>

{#if deletePopupOpen}
  <DeletePopup
    message={`Сигурни ли сте че искате да изтриете процедура: ${procedure.type}`}
    confirmAction={onDelete}
    onComplete={() => {
      deletePopupOpen = false;
    }}
  />
{/if}
