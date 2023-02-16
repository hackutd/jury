use rocket::{
    fairing::{Fairing, Info, Kind},
    http::Header,
    response::content::RawHtml,
    Request, Response,
};

pub struct CORS;

#[rocket::async_trait]
impl Fairing for CORS {
    fn info(&self) -> Info {
        Info {
            name: "Add CORS headers to responses",
            kind: Kind::Response,
        }
    }

    async fn on_response<'r>(&self, _request: &'r Request<'_>, response: &mut Response<'r>) {
        response.set_header(Header::new("Access-Control-Allow-Origin", "*"));
        response.set_header(Header::new(
            "Access-Control-Allow-Methods",
            "POST, GET, PATCH, OPTIONS",
        ));
        response.set_header(Header::new("Access-Control-Allow-Headers", "*"));
        response.set_header(Header::new("Access-Control-Allow-Credentials", "true"));
    }
}

/// Catches all OPTION requests in order to get the CORS related Fairing triggered.
#[rocket::options("/<_..>")]
pub fn all_options() {
    /* Intentionally left empty */
}

// Main client route
#[rocket::get("/<_..>", rank = 12)]
pub fn home() -> RawHtml<&'static str> {
    RawHtml(include_str!("../../public/index.html"))
}
