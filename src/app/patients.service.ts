import { Injectable } from '@angular/core';
import * as uuid from 'uuid';


export class Procedure {
  Id: string = "";
  Type: string = "";
  Date: Date = new Date();
  Details: String = "";
}

export class Patient {
  Id: string = ""
  Type: string = ""
  Name: string = ""
  Owner: string = ""
  Procedures: string[] = [];
  Gender: "male" | "female" | "unknown" = "unknown";
  BirthDate: Date = new Date();
  IdNumber: string = "0";
  LastModified: Date = new Date();
}

@Injectable({
  providedIn: 'root'
})
export class PatientsService {
  private patients = new Map<string, Patient>()
  private procedures = new Map<string, Procedure>()
  private types: String[] = ["Куче", "Котка", "Прица", "Rabbit"];

  private generateMockData() {
    this.procedures.set("0510ee92-b30e-4bff-a6d9-6af70b0e6acc", {
      Id: "0510ee92-b30e-4bff-a6d9-6af70b0e6acc",
      Type: 'Преглед',
      Date: new Date('2023-01-15'),
      Details: 'Routine check-up. All vitals are normal.'
    });
    this.procedures.set("b7306c26-8b48-4a8b-83c9-a2425c117364", {
      Id: "b7306c26-8b48-4a8b-83c9-a2425c117364",
      Type: 'Ваксина',
      Date: new Date('2023-02-20'),
      Details: 'Administered rabies vaccine.'
    });
    this.procedures.set("661a3ed6-c2a8-4b55-a3e2-d2f11afedbd1", {
      Id: "661a3ed6-c2a8-4b55-a3e2-d2f11afedbd1",
      Type: 'Кръвно изледване',
      Date: new Date('2023-03-10'),
      Details: 'Blood test for heartworm. Results are negative.'
    });

    let procedures = Array.from(this.procedures.keys());
    let patient: Patient = {
      Id: "1",
      Type: "Куче",
      Name: "Buddy",
      Owner: "John Doe",
      Procedures: procedures,
      Gender: "male",
      BirthDate: new Date("2018-01-15"),
      IdNumber: "482736194",
      LastModified: new Date("2023-01-01"),
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "2",
      Type: "Котка",
      Name: "Whiskers",
      Owner: "Jane Smith",
      Procedures: procedures,
      Gender: "female",
      BirthDate: new Date("2019-05-20"),
      IdNumber: "193847562",
      LastModified: new Date("2023-02-15")
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "3",
      Type: "Прица",
      Name: "Tweety",
      Owner: "Alice Johnson",
      Procedures: procedures,
      Gender: "unknown",
      BirthDate: new Date("2020-07-30"),
      IdNumber: "758392014",
      LastModified: new Date("2023-03-10")
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "4",
      Type: "Заек",
      Name: "Thumper",
      Owner: "Bob Brown",
      Procedures: procedures,
      Gender: "female",
      BirthDate: new Date("2021-11-05"),
      IdNumber: "620485731",
      LastModified: new Date("2023-04-20")
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "5",
      Type: "Куче",
      Name: "Max",
      Owner: "Charlie Davis",
      Procedures: procedures,
      Gender: "male",
      BirthDate: new Date("2017-03-10"),
      IdNumber: "374829105",
      LastModified: new Date("2023-05-25")
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "6",
      Type: "Куче",
      Name: "Бъди",
      Owner: "Иван Иванов",
      Procedures: procedures,
      Gender: "male",
      BirthDate: new Date("2018-01-15"),
      IdNumber: "918273645",
      LastModified: new Date("2023-01-01")
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "7",
      Type: "Котка",
      Name: "Мърка",
      Owner: "Мария Петрова",
      Procedures: procedures,
      Gender: "female",
      BirthDate: new Date("2019-05-20"),
      IdNumber: "506172839",
      LastModified: new Date("2023-02-15")
    };
    patient = {
      Id: "8",
      Type: "Прица",
      Name: "Чурулик",
      Owner: "Александър Георгиев",
      Procedures: procedures,
      Gender: "unknown",
      BirthDate: new Date("2020-07-30"),
      IdNumber: "284756193",
      LastModified: new Date("2023-03-10")
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "9",
      Type: "Rabbit",
      Name: "Тупър",
      Owner: "Борислав Димитров",
      Procedures: procedures,
      Gender: "female",
      BirthDate: new Date("2021-11-05"),
      IdNumber: "739182645",
      LastModified: new Date("2023-04-20")
    };
    this.patients.set(patient.Id, patient);
    patient = {
      Id: "10",
      Type: "Куче",
      Name: "Макс",
      Owner: "Георги Василев",
      Procedures: procedures,
      Gender: "male",
      BirthDate: new Date("2017-03-10"),
      IdNumber: "561093827",
      LastModified: new Date("2023-05-25")
    };
    this.patients.set(patient.Id, patient);
  }
  constructor() {
    this.generateMockData();
  }


  public getPatients() {
    return Array.from(this.patients.values())
  }

  public getPatient(key: string): Patient | undefined {
    const patient = this.patients.get(key);
    return patient ? { ...patient } : undefined;
  }

  public getProcedure(key: string): Procedure | undefined {
    const procedure = this.procedures.get(key);
    return procedure ? { ...procedure } : undefined;
  }

  public getProcedures(keys: string[]): Procedure[] {
    var procedures: Procedure[] = [];
    var service = this;
    keys.forEach(function(key) {
      const procedure = service.procedures.get(key);
      if (procedure) {
        procedures.push(procedure);
      }
    })

    return procedures;
  }

  public addPatient(patient: Patient) {
    patient.Id = uuid.v4();
    this.patients.set(patient.Id, patient);
  }

  public addProcedure(patientId: string, procedure: Procedure) {
    procedure.Id = uuid.v4();
    this.procedures.set(procedure.Id, procedure);
    let patient = this.patients.get(patientId);
    patient?.Procedures.unshift(procedure.Id);
  }

  public updatePatient(patient: Patient) {
    this.patients.set(patient.Id, patient);
  }


  public updateProcedure(procedure: Procedure) {
    this.procedures.set(procedure.Id, procedure);
  }

  public getTypes(): String[] {
    return this.types;
  }
}