use std::iter;

use actix_cors::Cors;
use actix_web::{get, post, App, HttpResponse, HttpServer, Responder};

use actix_web::middleware::Logger;
use env_logger::Env;
use rand::Rng;
use types::{Gender, Patient, Procedure};
use uuid::Uuid;

use crate::types::ToViewPatientList;

mod types;

#[get("/v1/patient-list")]
async fn hello() -> impl Responder {
    let mock_patients: Vec<Patient> = (1..=10).map(generate_mock_patient).collect();
    let patient = mock_patients.as_view();
    HttpResponse::Ok().body(serde_json::to_string_pretty(&patient).unwrap())
}

#[post("/v1/patient-list")]
async fn echo(req_body: String) -> impl Responder {
    HttpResponse::Ok().body(req_body)
}

fn generate_mock_procedure() -> Procedure {
    let mut rng = rand::thread_rng();
    Procedure {
        id: Uuid::new_v4(),
        r#type: "Checkup".to_string(),
        date: format!("2023-0{}-{}", rng.gen_range(1..10), rng.gen_range(10..29)),
        details: "General health check".to_string(),
    }
}

fn generate_mock_patient(index: usize) -> Patient {
    let mut rng = rand::thread_rng();
    Patient {
        id: Uuid::new_v4(),
        r#type: "Dog".to_string(),
        name: format!("Dog {}", index),
        gender: if index % 2 == 0 {
            Gender::Male
        } else {
            Gender::Female
        },
        birth_date: format!(
            "201{}-0{}-{}",
            rng.gen_range(0..10),
            rng.gen_range(1..10),
            rng.gen_range(10..29)
        ),
        chip_id: format!("100{}", rng.gen_range(100000000..999999999)),
        weight: rng.gen_range(5.0..30.0),
        castrated: index % 3 == 0,
        last_modified: "2023-04-01".to_string(),
        note: "Very friendly".to_string(),
        owner: format!("Owner {}", index),
        owner_phone: format!("+100000{}", index),
        procedures: iter::repeat_with(generate_mock_procedure)
            .take(rng.gen_range(1..4))
            .collect(),
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::init_from_env(Env::default().default_filter_or("info"));
    HttpServer::new(|| {
        App::new()
            .wrap(
                Cors::default()
                    .allow_any_origin()
                    .allow_any_method()
                    .allow_any_header(),
            )
            .wrap(Logger::default())
            .wrap(Logger::new("%a %{User-Agent}i"))
            .service(hello)
            .service(echo)
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
