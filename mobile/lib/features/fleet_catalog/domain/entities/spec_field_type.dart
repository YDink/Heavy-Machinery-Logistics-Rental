enum SpecFieldType {
  string,
  decimal,
  integer,
  boolean,
  enumeration;

  static SpecFieldType fromString(String val) {
    switch (val.toUpperCase()) {
      case 'DECIMAL':
        return SpecFieldType.decimal;
      case 'INTEGER':
        return SpecFieldType.integer;
      case 'BOOLEAN':
        return SpecFieldType.boolean;
      case 'ENUM':
        return SpecFieldType.enumeration;
      default:
        return SpecFieldType.string;
    }
  }
}
