import { type ViewPatient, type ViewProcedure, type ViewSetting } from "$lib/types";
import { writable } from "svelte/store";
import { dev } from '$app/environment';

var server = "/v1";
if (dev) {
  console.log("Running in development mode");
  server = "http://127.0.0.1:8001/v1";
}

var patients: Array<ViewPatient> | null = null;
var patientTypes: Array<ViewSetting> | null = null;
var procedureTypes: Array<ViewSetting> | null = null;
var patientFolders: Array<ViewSetting> | null = null;
export var patientsStore = writable(new Array<ViewPatient>())
export var settingsPatients = writable(new Array<ViewSetting>())
export var settingsProcedures = writable(new Array<ViewSetting>())
export var settingsFolders = writable(new Array<ViewSetting>())

export async function fetchPatients() {
  var response = await fetch(`${server}/patient-list`);
  patients = await response.json();
  patientsStore.set(patients!);
}

export async function fetchSettings() {
  var promise1 = fetch(`${server}/patient-types`);
  var promise2 = fetch(`${server}/procedure-types`);
  var promise3 = fetch(`${server}/patient-folder`);
  var [patientsResponse, proceduresResponse, foldersResponse] = await Promise.all([promise1, promise2, promise3])
  patientTypes = await patientsResponse.json();
  procedureTypes = await proceduresResponse.json();
  patientFolders = await foldersResponse.json();
  settingsPatients.set(patientTypes!)
  settingsProcedures.set(procedureTypes!)
  settingsFolders.set(patientFolders!)
}

export async function getPatient(id: Number): Promise<ViewPatient | undefined> {
  if (!patients) {
    await fetchPatients();
  }
  var patient = patients!.find((patient) => patient.id === id);

  return { ...patient };
}

export async function getProcedure(patientId: Number, id: Number): Promise<ViewProcedure | undefined> {
  if (!patients) {
    await fetchPatients();
  }
  var procedure = patients!.find((patient) => patient.id === patientId)?.procedures?.find((procedure) => procedure.id == id)
  if (procedure) {
    return { ...procedure };
  }

  return undefined;
}

// Setters

export async function updatePatient(patient: ViewPatient) {
  const index = patients!.findIndex((p) => p.id === patient.id);
  patients![index] = patient;
  patientsStore.set(patients!);
  console.log("update patient", patient);

  await fetch(`${server}/patient`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patient)
  });
}

export async function addPatient(patient: ViewPatient) : Promise<ViewPatient> {
  // TODO: add proper id
  // patient.id = patients?.length! + Math.floor(Math.random() * 1000);

  var response = await fetch(`${server}/patient`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patient)
  });
  var newPatient = await response.json();
  patients?.push(newPatient);
  patientsStore.set(patients!);
  console.log("add patient", newPatient);
  return newPatient;
}

export async function deletePatient(id: Number) {
  var index = patients!.findIndex((patient) => patient.id == id);
  patients!.splice(index, 1);
  patientsStore.set(patients!);
  console.log("delete patient", id);

  await fetch(`${server}/patient/${id}`, {
    method: "DELETE"
  });
}

export async function addProcedure(patientId: number, procedure: ViewProcedure) {
  var patient = patients!.find((patient) => patient.id === patientId);
  // TODO: add proper id
  // procedure.id = patient?.procedures?.length! + Math.floor(Math.random() * 1000);
  patient?.procedures?.push(procedure);
  patientsStore.set(patients!);
  procedure.patientId = patientId;
  console.log("add procedure", procedure);
  
  await fetch(`${server}/procedure/`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(procedure)
  });
}

export async function updateProcedure(patientId: number, procedure: ViewProcedure) {
  var patient = patients!.find((patient) => patient.id === patientId);
  const index = patient?.procedures!.findIndex((p) => p.id === procedure.id);
  patient!.procedures![index!] = procedure;
  patientsStore.set(patients!);
  procedure.patientId = patientId;
  console.log("update procedure", procedure);
  
  await fetch(`${server}/procedure/`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(procedure)
  });
}

export async function deleteProcedure(patientId: Number, id: Number) {
  var patient = patients!.find((patient) => patient.id === patientId);
  var index = patient?.procedures!.findIndex((procedure) => procedure.id == id);
  patient?.procedures!.splice(index!, 1);
  patientsStore.set(patients!);
  console.log("delete procedure", id);
  
  await fetch(`${server}/procedure/${id}`, {
    method: "DELETE"
  });
}

export async function updatePatientSettings(patientTypes: ViewSetting[]) {
  console.log("update patient type", patientTypes);
  
  await fetch(`${server}/settings`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patientTypes)
  });
}

export async function updateProcedureSettings(procedureTypes: ViewSetting[]) {
  console.log("update procedure type", procedureTypes);
  
  await fetch(`${server}/settings`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(procedureTypes)
  });
}

export async function updatePatientFolders(patientFolders: ViewSetting[]) {
  console.log("patient folders", patientFolders);
  
  await fetch(`${server}/settings`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patientFolders)
  });
}

export async function deleteSetting(id: Number) {
  console.log("deleting setting", id);
  
  await fetch(`${server}/setting/${id}`, {
    method: "DELETE",
  });
}

fetchPatients();
fetchSettings();
