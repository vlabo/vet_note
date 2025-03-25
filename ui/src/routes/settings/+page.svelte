<script lang="ts">
  import BackButton from "$lib/BackButton.svelte";
  import DeletePopup from "$lib/DeletePopup.svelte";
  import List from "./components/list.svelte";
  import {
    settingsPatients,
    settingsProcedures,
    updatePatientSettings,
    updateProcedureSettings,
  } from "$lib/DataService";

  // Style
  var patientTypes = settingsPatients;
  var procedureTypes = settingsProcedures;
  var deletePopupOpen = false;
</script>

<div class="bg-gray-50">
  <header class="max-w-7xl mx-auto bg-white shadow flex items-center p-2">
    <BackButton />
    <h1 class="font-semibold mx-2">Настройки</h1>
  </header>

  <main class="max-w-7xl mx-auto px-4 py-2">
    <!-- Patient Types Section -->
    <List
      title="Тип животно"
      placeholder="Тип животно"
      settingType="PatientType"
      bind:items={patientTypes}
      onUpdate={updatePatientSettings}
    />

    <!-- Procedure Types Section -->
    <List
      title="Тип процедура"
      placeholder="Тип процедура"
      settingType="ProcedureType"
      bind:items={procedureTypes}
      onUpdate={updateProcedureSettings}
    />

    <div class="mt-4 mb-6">
      <!-- User Section -->
      <section>
        <h2 class="text-gray-700 font-medium mb-4">Потребител</h2>
        <button
          class="w-full bg-red-500 text-white py-2 rounded focus:outline-none"
        >
          Излез
        </button>
      </section>
    </div>
  </main>
</div>

{#if deletePopupOpen}
  <DeletePopup
    onComplete={() => {
      deletePopupOpen = false;
    }}
  />
{/if}
