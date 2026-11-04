enum PhotoAngle {
  front,
  rear,
  left,
  right,
  cabin,
  meter;

  String get code => name.toUpperCase();

  static PhotoAngle? fromCode(String code) {
    for (final a in PhotoAngle.values) {
      if (a.code == code.toUpperCase()) return a;
    }
    return null;
  }
}
