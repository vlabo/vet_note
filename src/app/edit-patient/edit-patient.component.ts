
import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { IonicModule } from '@ionic/angular';
import { Patient, PatientsService } from '../patients.service';
import { Location } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-edit-patient',
  templateUrl: './edit-patient.component.html',
  styleUrls: ['./edit-patient.component.scss'],
  standalone: true,
  imports: [CommonModule, FormsModule, IonicModule]
})
export class EditPatientComponent implements OnInit {
  patient: Patient = new Patient();
  types: String[] = [];
  newMode = false;

  constructor(
    private patientsService: PatientsService,
    private location: Location,
    private route: ActivatedRoute,
    private router: Router,
  ) { }

  ngOnInit(): void {
    this.route.url.subscribe(urlSegments => {
      var lastPathSegment = urlSegments[urlSegments.length - 1].path
      if(lastPathSegment === "new") {
        this.newMode = true;
      }
    });

    this.route.paramMap.subscribe(params => {
      let id = params.get('id');
      if (id === null) {
        return;
      }
      var patient = this.patientsService.getPatient(id);
      if (patient) {
        console.log(patient);
        this.patient = patient;
      }
    });
    this.types = this.patientsService.getTypes();
  }

  save(): void {
    if (this.newMode) {
      this.patientsService.addPatient(this.patient);
      this.newMode = false;
      this.router.navigate(["/patient", this.patient.Id], { replaceUrl: true })
    } else {
      this.patientsService.updatePatient(this.patient);
      this.location.back();
    }
  }

  cancel(): void {
    this.location.back();
  }
}