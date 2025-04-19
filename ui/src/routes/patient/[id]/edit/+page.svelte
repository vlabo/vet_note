<script lang="ts">
  import { type ViewPatient, type ViewSetting } from "$lib/types";
  import { page } from "$app/state";
  import { onMount } from "svelte";
  import {
    addPatient,
    getPatient,
    settingsPatients,
    updatePatient,
    deletePatient,
  } from "$lib/DataService";
  import { goto } from "$app/navigation";
  import type { Writable } from "svelte/store";
  import Icon from "$lib/icon.svelte";
  import { faTrash } from "@fortawesome/free-solid-svg-icons";
  import BackButton from "$lib/BackButton.svelte";
  import DeletePopup from "$lib/DeletePopup.svelte";

  console.log(page.params);
  var types: Writable<Array<ViewSetting>> = settingsPatients;
  var patient: ViewPatient = {
    id: undefined,
    type: null,
    name: null,
    gender: "unknown",
    age: null,
    chipId: null,
    weight: null,
    castrated: null,
    folder: null,
    indexFolder: null,
    note: null,
    owner: null,
    ownerPhone: null,
    procedures: [],
  };
  var newMode = false;
  var deletePopupOpen = false;
  var castrated = false;

  onMount(() => {
    if (page.params.id === "new") {
      newMode = true;
      types.subscribe(($items) => {
        if (patient.type === null) {
          patient.type = $items[0].value;
        }
      });
    } else {
      getPatient(Number(page.params.id)).then((p) => {
        if (p) {
          patient = p;
          castrated = p.castrated === 1;
        }
      });
    }
  });

  function update() {
    patient.castrated = castrated ? 1 : 0;
    if (newMode) {
      addPatient(patient).then((patient) => {
        goto(`/patient/${patient.id}`, { replaceState: true });
      });
      // goto(`/patient/${patient.id}`, { replaceState: true });
    } else {
      updatePatient(patient);
      history.back();
    }
  }

  function onDelete() {
    deletePatient(patient.id!);
    history.go(-2);
  }
</script>

<div class="min-h-screen">
  <!-- HEADER -->
  <div
    class="bg-white flex shadow p-4 max-w-7xl mx-auto justify-between items-center"
  >
    <BackButton />
    <h1
      class="absolute left-1/2 transform -translate-x-1/2 text-lg font-semibold"
    >
      {newMode ? "Нов Пациент" : "Редактирай Пациент"}
    </h1>
    {#if !newMode}
      <button
        on:click={() => (deletePopupOpen = true)}
        class="hover:bg-blue-200 text-red-500 text-sm font-medium px-5 py-3 rounded"
      >
        <Icon icon={faTrash} />
      </button>
    {/if}
  </div>

  <!-- FORM -->
  <div class="bg-white shadow p-4 max-w-7xl mx-auto items-center">
    <!-- Тип -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Тип:
        <select
          bind:value={patient.type}
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        >
          <!--TODO: add key :key="type.idx" -->
          {#each $types as type: ViewSetting}
            <option>{type.value}</option>
          {/each}
        </select>
      </label>
    </div>

    <!-- Име -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Име:
        <input
          bind:value={patient.name}
          type="text"
          placeholder="Име"
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        />
      </label>
    </div>

    <!-- Пол -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Пол:
        <select
          bind:value={patient.gender}
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        >
          <option value="unknown">Неопределен</option>
          <option value="male">Мъжки</option>
          <option value="female">Женски</option>
        </select>
      </label>
    </div>

    <!-- Тегло (кг) -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Тегло (кг):
        <input
          bind:value={patient.weight}
          type="number"
          inputmode="numeric"
          placeholder="0"
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        />
      </label>
    </div>

    <!-- Кастрирано -->
    <div class="flex items-center">
      <input
        bind:checked={castrated}
        type="checkbox"
        id="castrated"
        class="h-4 w-4 text-blue-600 border rounded m-3"
      />
      <label for="castrated" class="ml-2 block text-sm text-gray-700"
        >Кастрирано</label
      >
    </div>

    <!-- Възраст (години) -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Възраст (години):
        <input
          bind:value={patient.age}
          type="number"
          inputmode="numeric"
          placeholder="0"
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        />
      </label>
    </div>

    <!-- Име на собственик -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Име на собственик:
        <input
          bind:value={patient.owner}
          type="text"
          placeholder="Име"
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        />
      </label>
    </div>

    <!-- Телефон на собственик -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Телефон на собственик:
        <input
          bind:value={patient.ownerPhone}
          type="text"
          placeholder="0812345679"
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        />
      </label>
    </div>

    <!-- Чип -->
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1"
        >Чип:
        <input
          bind:value={patient.chipId}
          type="text"
          placeholder="100..."
          class="w-full border rounded shadow-sm focus:ring-blue-500 focus:border-blue-500 p-3"
        />
      </label>
    </div>
  </div>

  <!-- SAVE BUTTON -->
  <div class="p-4 max-w-7xl mx-auto items-center">
    <button
      class="w-full bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded"
      on:click={update}
    >
      Запиши
    </button>
  </div>
</div>

{#if deletePopupOpen}
  <DeletePopup
    message={`Сигурни ли сте че искате да изтриете пациента ${patient.name}`}
    confirmAction={onDelete}
    onComplete={() => {
      deletePopupOpen = false;
    }}
  />
{/if}
