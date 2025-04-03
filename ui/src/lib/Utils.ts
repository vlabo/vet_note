import type { ViewPatient } from "./types";

export function formatDate(date: Date|string|undefined) : string {
  if(!date) {
    return "";
  }
  if(typeof date == "string") {
     date = new Date(date);   
  } 
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0"); // getMonth() returns 0-based month
  const day = String(date.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}

export interface ViewPatientHighlighted extends ViewPatient {
  highlightedFields: any;
}
