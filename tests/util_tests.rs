use jury::util::crowd_bt;

#[test]
fn math_test() {
    let (a, b, c, d, e, f) = crowd_bt::update(20.0, 12.2, 4.2, 1.0, 3.11, 0.65);
    assert_eq!(a, 20.29342378562617);
    assert_eq!(b, 12.144888310192417);
    assert_eq!(c, 4.28143039999674);
    assert_eq!(d, 0.9529174440716865);
    assert_eq!(e, 3.057070240002119);
    assert_eq!(f, 0.6301076201202875);
}
