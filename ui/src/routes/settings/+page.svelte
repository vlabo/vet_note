<script lang="ts">
  import BackButton from "$lib/BackButton.svelte";
  import DeletePopup from "$lib/DeletePopup.svelte";
  import List from "./components/list.svelte";
  import {
    settingsFolders,
    settingsPatients,
    settingsProcedures,
    updatePatientSettings,
    updateProcedureSettings,
    updatePatientFolders,
  } from "$lib/DataService";

  // Style
  var patientTypes = settingsPatients;
  var procedureTypes = settingsProcedures;
  var patientFolders = settingsFolders;
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

    <!-- Patient folders -->
    <List
      title="Папки на пациенти"
      placeholder="име на папка"
      settingType="PatientFolder"
      bind:items={patientFolders}
      onUpdate={updatePatientFolders}
    />
  </main>
</div>

{#if deletePopupOpen}
  <DeletePopup
    onComplete={() => {
      deletePopupOpen = false;
    }}
  />
{/if}
