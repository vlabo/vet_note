import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AlertController, IonicModule } from '@ionic/angular';
import { CommonModule, Location } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PatientsService } from '../patients.service';

import { addIcons } from "ionicons";
import { arrowBack } from "ionicons/icons";
import { ViewProcedure } from '../types';

@Component({
  selector: 'app-procedure',
  templateUrl: './procedure.component.html',
  styleUrls: ['./procedure.component.scss'],
  imports: [IonicModule, CommonModule, FormsModule],
  standalone: true,
})
export class ProcedureComponent implements OnInit {
  procedure: ViewProcedure;
  isEditMode: boolean = false;
  isNewMode: boolean = false;

  procedureTypes: string[] = [];
  patientId: string | null = null;
  date: string = "";

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private location: Location,
    private patientsService: PatientsService,
    private alertController: AlertController,
  ) {
    addIcons({ "arrow-back": arrowBack })

    this.procedure = {
      id: "",
      type: "",
      date: "",
      details: "",
      patientId: "",
    };
  }

  ngOnInit(): void {

    this.patientsService.getProcedureTypes().subscribe({
      next: types => {
        if(types) {
          this.procedureTypes = types;
        }
      }
    });
    
    this.route.paramMap.subscribe(async paramMap => {
      const procedureId = paramMap.get('procedureId');
      if (procedureId) {
        this.patientsService.getProcedure(procedureId).subscribe({
          next: procedure => {
            this.procedure = procedure;
            console.log(this.procedure);
            this.date = this.procedure.date;
            this.patientId = procedure.patientId;
          }
        })
      } else {
        this.patientId = paramMap.get('patientId');
      }
    });

    this.route.queryParamMap.subscribe(queryParams => {
      if (!this.isNewMode) {
        this.isEditMode = queryParams.get('edit') === 'true';
      }
    });

    this.route.data.subscribe(data => {
      console.log(data);
      this.isNewMode = data["newMode"];
      if (this.isNewMode) {
        this.isEditMode = true;
      }
    });
  }

  goBack(): void {
    this.location.back();
  }

  enableEditMode() {
    this.isEditMode = true;
    this.router.navigate([], { queryParams: { edit: true }, queryParamsHandling: 'merge' });
  }

  saveProcedure(): void {
    this.patientsService.updateProcedure(this.patientId!, this.procedure!).subscribe({
      next: procedure => {
        this.procedure = procedure;
        if (this.isNewMode) {
          this.router.navigate(["procedure", this.procedure!.id], { queryParams: { edit: false }, replaceUrl: true });
        } else if (this.isEditMode) {
          this.goBack();
        }
      }
    })
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
            this.patientsService.deleteProcedure(this.procedure.id).subscribe({
              next: _ => {
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