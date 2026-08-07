/// Type-safe model for the `products` entity.
class Product {
  final String? id;
  final String? tenantId;
  final String? sku;
  final String name;
  final String? description;
  final String? categoryId;
  final String? unitId;
  final double costPrice;
  final double salePrice;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  const Product({
    this.id,
    this.tenantId,
    this.sku,
    required this.name,
    this.description,
    this.categoryId,
    this.unitId,
    this.costPrice = 0,
    this.salePrice = 0,
    this.createdAt,
    this.updatedAt,
  });

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id']?.toString(),
      tenantId: json['tenant_id']?.toString(),
      sku: json['sku']?.toString(),
      name: json['name']?.toString() ?? '',
      description: json['description']?.toString(),
      categoryId: json['category_id']?.toString(),
      unitId: json['unit_id']?.toString(),
      costPrice: _toDouble(json['cost_price']),
      salePrice: _toDouble(json['sale_price']),
      createdAt: _parseDate(json['created_at']),
      updatedAt: _parseDate(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() => {
        if (id != null) 'id': id,
        if (tenantId != null) 'tenant_id': tenantId,
        if (sku != null) 'sku': sku,
        'name': name,
        if (description != null) 'description': description,
        if (categoryId != null) 'category_id': categoryId,
        if (unitId != null) 'unit_id': unitId,
        'cost_price': costPrice,
        'sale_price': salePrice,
        if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
        if (updatedAt != null) 'updated_at': updatedAt!.toIso8601String(),
      };

  Product copyWith({
    String? id,
    String? tenantId,
    String? sku,
    String? name,
    String? description,
    String? categoryId,
    String? unitId,
    double? costPrice,
    double? salePrice,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Product(
      id: id ?? this.id,
      tenantId: tenantId ?? this.tenantId,
      sku: sku ?? this.sku,
      name: name ?? this.name,
      description: description ?? this.description,
      categoryId: categoryId ?? this.categoryId,
      unitId: unitId ?? this.unitId,
      costPrice: costPrice ?? this.costPrice,
      salePrice: salePrice ?? this.salePrice,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

double _toDouble(dynamic value) {
  if (value == null) return 0;
  if (value is double) return value;
  if (value is int) return value.toDouble();
  return double.tryParse(value.toString()) ?? 0;
}

DateTime? _parseDate(dynamic value) {
  if (value == null) return null;
  if (value is DateTime) return value;
  final str = value.toString();
  if (str.isEmpty) return null;
  return DateTime.tryParse(str);
}
