use serde::Deserialize;

#[derive(Deserialize)]
#[serde()]
pub struct Login<'r> {
    pub code: &'r str,
}

#[derive(Deserialize)]
#[serde()]
pub struct NewJudge<'r> {
    pub name: &'r str,
    pub email: &'r str,
    pub notes: &'r str,
}
