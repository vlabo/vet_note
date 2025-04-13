import { type ViewPatient, type ViewProcedure, type ViewSetting } from "$lib/types";
import { derived, writable } from "svelte/store";
import { dev } from '$app/environment';

var server = "/v1";
if (dev) {
  console.log("Running in development mode");
  server = "http://127.0.0.1:8001/v1";
}


export interface Folder {
  setting: ViewSetting;
  patients: Array<ViewPatient>;
}

var patients: Map<number, ViewPatient> = new Map();
var patientTypes: Array<ViewSetting> | null = null;
var procedureTypes: Array<ViewSetting> | null = null;
var patientFolders: Array<ViewSetting> | null = null;
export var patientsStore = writable(new Array<ViewPatient>())
export var settingsPatients = writable(new Array<ViewSetting>())
export var settingsProcedures = writable(new Array<ViewSetting>())
export var settingsFolders = writable(new Array<ViewSetting>())


function updatePatientStore(patients: Array<ViewPatient>) {
  patients.sort((a, b) => a.indexFolder! - b.indexFolder!);
  patientsStore.set(patients);
}

export var folders = derived([patientsStore, settingsFolders], ([$patientsStore, $settingsFolders]) => {
  return $settingsFolders.map((folder) => {
    const patients = $patientsStore.filter((patient) => patient.folder == folder.id);
    return { setting: folder, patients };
  });
});

export async function fetchPatients() {
  let response = await fetch(`${server}/patient-list`);
  let patientsArray = await response.json();
  patientsArray.forEach((p: ViewPatient) => {
    patients.set(p.id, p);
  })
  updatePatientStore(patients!.values().toArray());
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

export async function getPatient(id: number): Promise<ViewPatient | undefined> {
  if (patients.size == 0) {
    await fetchPatients();
  }
  var patient = patients!.get(id); 
  if (!patient) return undefined;

  return { ...patient };
}

export async function getProcedure(patientId: number, id: number): Promise<ViewProcedure | undefined> {
  if (!patients) {
    await fetchPatients();
  }
  var procedure = patients!.get(patientId)?.procedures.find((procedure: ViewProcedure) => procedure.id == id); 
  if (procedure) {
    return { ...procedure };
  }

  return undefined;
}

// Setters

export async function updatePatient(patient: ViewPatient) {
  patients!.set(patient.id, patient);
  updatePatientStore(patients!.values().toArray());
  console.log("update patient", patient);

  await fetch(`${server}/patient`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patient)
  });
}

export async function updatePatients(updated: Array<ViewPatient>) {
  updated.forEach((patient => {
    patients!.set(patient.id, patient);
  }));
  updatePatientStore(patients!.values().toArray());
  console.log("update patient", updated);

  await fetch(`${server}/patients/`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(updated)
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
  patients?.set(newPatient.id, newPatient);
  updatePatientStore(patients!.values().toArray());
  console.log("add patient", newPatient);
  return newPatient;
}

export async function deletePatient(id: number) {
  patients?.delete(id);
  updatePatientStore(patients!.values().toArray());
  console.log("delete patient", id);

  await fetch(`${server}/patient/${id}`, {
    method: "DELETE"
  });
}

export async function addProcedure(patientId: number, procedure: ViewProcedure) {
  // TODO: add proper id
  // procedure.id = patient?.procedures?.length! + Math.floor(Math.random() * 1000);
  patients!.get(patientId)!.procedures?.push(procedure);
  updatePatientStore(patients!.values().toArray());
  procedure.patientId = patientId;
  console.log("add procedure", procedure);
  
  await fetch(`${server}/procedure/`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(procedure)
  });
}

export async function updateProcedure(patientId: number, procedure: ViewProcedure) {
  var patient = patients!.get(patientId);
  const index = patient?.procedures!.findIndex((p) => p.id === procedure.id);
  patient!.procedures![index!] = procedure;
  updatePatientStore(patients!.values().toArray());
  procedure.patientId = patientId;
  console.log("update procedure", procedure);
  
  await fetch(`${server}/procedure/`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(procedure)
  });
}

export async function deleteProcedure(patientId: number, id: number) {
  var patient = patients!.get(patientId);
  var index = patient?.procedures!.findIndex((procedure) => procedure.id == id);
  patient?.procedures!.splice(index!, 1);
  updatePatientStore(patients!.values().toArray());
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
