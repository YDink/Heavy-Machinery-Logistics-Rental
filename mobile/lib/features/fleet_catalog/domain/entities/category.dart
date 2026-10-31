import 'package:equatable/equatable.dart';

class Category extends Equatable {
  final String id;
  final String code;
  final String name;
  final String iconUrl;
  final String? parentId;

  const Category({
    required this.id,
    required this.code,
    required this.name,
    required this.iconUrl,
    this.parentId,
  });

  @override
  List<Object?> get props => [id, code, name, iconUrl, parentId];
}
