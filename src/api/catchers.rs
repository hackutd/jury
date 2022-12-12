use rocket::{catch, response::Redirect, Request};

#[catch(401)]
pub fn unauthorized(req: &Request) -> Redirect {
    let route = match req.route() {
        Some(r) => r,
        None => {
            return Redirect::to("/");
        }
    };
    let path = route.uri.path();
    if path.contains("judge") {
        Redirect::to("/judge/login")
    } else {
        Redirect::to("/admin/login")
    }
}
