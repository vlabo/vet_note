<script lang="ts">
  import Icon from "$lib/icon.svelte";
  import {
    faCalendarDays,
    faMars,
    faVenus,
    faWeightHanging,
    faCheck,
    faMicrochip,
    faUser,
    faPhone,
    faAngleRight,
    faEdit,
    faPlus,
  } from "@fortawesome/free-solid-svg-icons";
  import { type ViewPatient } from "$lib/types";
  import BackButton from "$lib/BackButton.svelte";
  import { page } from "$app/state";
  import { onMount } from "svelte";
  import { getPatient, updatePatient } from "$lib/DataService";
  import { formatDate } from "$lib/Date";

  console.log(page.params);

  var patient: ViewPatient = $state({});

  onMount(() => {
    getPatient(Number(page.params.id)).then((p) => {
      if (p) {
        patient = p;
      }
    });
  });

  let canUpdateNote = true;

  function onNoteChange() {
    if (!canUpdateNote) return;

    canUpdateNote = false;
    setTimeout(() => {
      canUpdateNote = true;
    }, 500);

    updatePatient(patient);
  }
</script>

<div class="min-h-screen">
  <!-- Only show if patient exists -->
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <header class="bg-white shadow p-4">
      <div class="flex items-center justify-between">
        <BackButton />
        <div class="text-center">
          <h1 class="text-xl font-bold">{patient.name}</h1>
          <p class="text-sm text-gray-600">{patient.type}</p>
        </div>
        <div class="space-x-2">
          <a
            class="text-black focus:outline-none mx-1 px-5 py-4 hover:bg-blue-100 hover:rounded-lg"
            href={`/patient/${patient.id}/edit`}
          >
            <Icon icon={faEdit} />
          </a>
        </div>
      </div>
    </header>

    <!-- Patient details -->
    <div class="flex flex-wrap items-center gap-x-8 gap-y-4 p-4">
      <!-- Age -->
      {#if patient.age}
        <div class="flex items-center">
          <Icon icon={faCalendarDays} />
          <span class="ml-4">{patient.age}г.</span>
        </div>
      {/if}

      <!-- Gender -->
      <div class="flex items-center">
        {#if patient.gender === "male"}
          <div class="flex items-center">
            <Icon icon={faMars} />
            <span class="ml-2">Мъжки</span>
          </div>
        {:else if patient.gender === "female"}
          <div class="flex items-center">
            <Icon icon={faVenus} />
            <span class="ml-2">Женски</span>
          </div>
        {/if}
      </div>

      <!-- Weight -->
      {#if patient.weight}
        <div class="flex items-center">
          <Icon icon={faWeightHanging} />
          <span class="ml-2">{patient.weight}кг</span>
        </div>
      {/if}

      <!-- Castration status -->
      {#if patient.castrated == 1}
        <div class="flex items-center">
          <div class="flex items-center">
            <Icon icon={faCheck} />
            <span class="ml-2">Кастрирано</span>
          </div>
        </div>
      {/if}

      <!-- Chip ID -->
      {#if patient.chipId}
        <div class="flex items-center">
          <Icon icon={faMicrochip} />
          <span class="ml-2">{patient.chipId}</span>
        </div>
      {/if}

      <!-- Owner -->
      {#if patient.owner}
        <div class="flex items-center">
          <Icon icon={faUser} />
          <span class="ml-2">{patient.owner}</span>
        </div>
      {/if}

      <!-- Owner Phone -->
      {#if patient.ownerPhone}
        <div class="flex items-center">
          <Icon icon={faPhone} />
          <span class="ml-2">{patient.ownerPhone}</span>
        </div>
      {/if}
    </div>
    <div class="p-4 space-y-4">
      <!-- Note Textarea -->
      <div>
        <textarea
          bind:value={patient.note}
          onchange={onNoteChange}
          placeholder="Бележка"
          class="w-full p-2 border rounded focus:outline-none focus:ring resize-none"
        ></textarea>
      </div>

      <!-- Procedures History -->
      <div class="mt-6">
        <h2 class="text-lg font-semibold mb-4">История</h2>

        <div class="space-y-2">
          <a
            class="w-full border-2 border-dashed border-gray-300 p-2 rounded flex justify-center items-center focus:outline-none"
            href={`/patient/${patient.id}/procedure/new`}
          >
            <Icon icon={faPlus} class="mx-2" />
          </a>

          {#each patient.procedures! as procedure}
            <a
              class="flex items-center p-3 bg-white rounded border cursor-pointer hover:bg-gray-50"
              href={`/patient/${patient.id}/procedure/${procedure.id}`}
            >
              <div class="flex items-center space-x-4">
                <div>
                  <h3 class="font-semibold">{procedure.type}</h3>
                  <p class="text-sm text-gray-600"></p>
                </div>
                <div class="text-gray-500">{formatDate(procedure.date)}</div>
              </div>
              <Icon icon={faAngleRight} class="text-black ml-auto" />
            </a>
          {/each}
        </div>
      </div>
    </div>
  </div>
</div>
