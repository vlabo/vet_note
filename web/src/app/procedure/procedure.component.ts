import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AlertController } from '@ionic/angular';
import * as common from '@angular/common';
import { PatientsService } from '../patients.service';

import { addIcons } from "ionicons";
import { arrowBack } from "ionicons/icons";
import { ViewProcedure } from '../types';

@Component({
  selector: 'app-procedure',
  templateUrl: './procedure.component.html',
  styleUrls: ['./procedure.component.scss'],
})
export class ProcedureComponent implements OnInit {
  originalProcedure: ViewProcedure;
  procedure: ViewProcedure;
  mode: "new" | "edit" | "view" = "view";

  procedureTypes: string[] = [];
  date: string = "";

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private location: common.Location,
    private patientsService: PatientsService,
    private alertController: AlertController,
  ) {
    addIcons({ "arrow-back": arrowBack })

    this.originalProcedure = {
      patientId: 0,
    };
    this.procedure = {
      date: new Date().toISOString(),
      patientId: 0,
    };
  }

  ngOnInit(): void {

    this.patientsService.getProcedureTypes().subscribe({
      next: types => {
        if(types) {
          this.procedureTypes = types.map((type) : string => { return type.value; });
        }
      }
    });
    
    this.route.paramMap.subscribe(async paramMap => {
      const procedureId = paramMap.get('procedureId');
      if (procedureId) {
        this.patientsService.getProcedure(procedureId).subscribe({
          next: procedure => {
            this.originalProcedure = procedure;
            this.procedure = {...procedure};
            this.date = this.procedure.date!;
          }
        })
      } else {
        this.procedure.patientId = Number(paramMap.get('patientId'));
      }
    });

    this.route.queryParamMap.subscribe(queryParams => {
        if(queryParams.get('edit') === 'true') {
          this.mode = "edit";
        } else {
          this.mode = "view";
          
        }
    });

    this.route.data.subscribe(data => {
      if(data["newMode"]) {
        this.mode = "new";
      }
    });
  }

  goBack(): void {
    this.location.back();
  }

  enableEditMode() {
    this.router.navigate([], { queryParams: { edit: true }, queryParamsHandling: 'merge' });
  }

  saveProcedure(): void {
    this.patientsService.updateProcedure(this.getChangedValues(this.originalProcedure, this.procedure)).subscribe({
      next: procedure => {
        this.procedure = procedure;
        if (this.mode == "new") {
          console.log(this.procedure);
          this.router.navigate(["procedure", this.procedure!.id], { queryParams: { edit: false }, replaceUrl: true });
        } else if (this.mode == "edit") {
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
            this.patientsService.deleteProcedure(this.procedure.id!).subscribe({
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

  getChangedValues(original: ViewProcedure, updated: ViewProcedure): ViewProcedure {
    const changes: ViewProcedure = {
      id: original.id,
      patientId: original.patientId,
    };

    for (const key in updated) {
      // @ts-ignore
      if (original[key] !== updated[key]) {
        // @ts-ignore
        changes[key] = updated[key];
      }
    }

    return changes;
  }
}