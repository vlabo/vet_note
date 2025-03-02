import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { PatientsService } from '../patients.service';
import { AlertController, IonModal } from '@ionic/angular';

import { faCakeCandles, faUser, faCalendarDays, faMarsAndVenus, faMicrochip, faMars, faVenus, faPhone, faWeightHanging, faCheck, faX } from '@fortawesome/free-solid-svg-icons';

import { add, create, checkmark, close, arrowBack, chevronForward, paw, person, calendar, heartHalf } from 'ionicons/icons';
import { addIcons } from "ionicons";
import { Location } from '@angular/common';
import { ViewPatient, ViewProcedure } from '../types';

@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss'],
})
export class PatientComponent implements OnInit, OnDestroy {

  // Icons
  cake = faCakeCandles;
  user = faUser;
  calendar = faCalendarDays;
  gender = faMarsAndVenus;
  microchip = faMicrochip;
  mars = faMars;
  venus = faVenus;
  phone = faPhone;
  weight = faWeightHanging;
  check = faCheck;
  xIcon = faX;

  patient: ViewPatient | undefined = undefined;
  originalNote: string = "";
  procedures: ViewProcedure[] = [];
  isViewProcedure = false;

  selectedProcedure: ViewProcedure | null = null;

  @ViewChild('procedureDetailsModal', { static: true }) modal!: IonModal;

  constructor(private route: ActivatedRoute,
    private router: Router,
    private location: Location,
    private patientService: PatientsService,
    private alertController: AlertController,
  ) {
    addIcons({
      "add": add,
      "create": create,
      "checkmark": checkmark,
      "close": close,
      "arrow-back": arrowBack,
      "chevron-forward": chevronForward,
      "paw": paw,
      "person": person,
      "calendar": calendar,
      "heart-half": heartHalf,
    });
  }

  ngOnInit() {
    this.route.paramMap.subscribe(async params => {
      let id = params.get('id');
      if (id === null) {
        return;
      }
      this.patientService.getPatient(id).subscribe({
        next: patient => {
          this.patient = patient;
          this.originalNote = patient.note || "";
        }
      })
    });
  }

  ngOnDestroy() {
    this.isViewProcedure = false;
    this.updateNote();
  }

  openEdit() {
    this.router.navigate(["/patient", this.patient?.id, "edit"]);
  }

  goBack(): void {
    this.location.back();
  }

  setOpen(isOpen: boolean) {
    this.isViewProcedure = isOpen;
  }

  openViewProcedure(procedure: ViewProcedure) {
    if (this.patient) {
      this.router.navigate(['/procedure', procedure.id]);
    }
  }

  createProcedure(): void {
    if (this.patient) {
      this.router.navigate(['/procedure', this.patient.id, 'new']);
      this.isViewProcedure = false;
    }
  }

  updateNote() {
    if (this.patient && this.patient.note !== this.originalNote) {
      let patient: ViewPatient = {
        id: this.patient.id,
        note: this.patient.note,
      }
      this.patientService.updatePatient(patient).subscribe({
        next: _ => {
          this.patientService.triggerPatientListReload();
        }
      });
      this.originalNote = this.patient.note || "";
    }
  }

  isValidDate(dateString: string | undefined): boolean {
    if (!dateString) {
      return false;
    }
    const date = new Date(dateString);
    return !isNaN(date.getTime());
  }

  async presentDeleteConfirm() {
    const alert = await this.alertController.create({
      header: 'Потвърди изтриване',
      message: 'Сигурни ли сте, че искате да изтриете този запис?',
      buttons: [
        {
          text: 'Отказ',
          role: 'cancel',
          cssClass: 'secondary',
          handler: () => { }
        },
        {
          text: 'Изтрий',
          handler: () => {
            this.patientService.deletePatient(this.patient!.id!).subscribe({
              next: _ => {
                this.patient = undefined;
                this.goBack();
              }
            })
          }
        }
      ]
    });

    await alert.present();
  }

}
