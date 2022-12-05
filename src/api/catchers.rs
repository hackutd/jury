use rocket::{catch, response::Redirect, Request};

#[catch(401)]
pub fn unauthorized(req: &Request) -> Redirect {
    let route = match req.route() {
        Some(r) => r,
        None => {
            return Redirect::to("/");
        }
    };
    Redirect::to(format!("{}/login", route.uri.path()))
}
