/// Type-safe model for the `payment_methods` entity.
class PaymentMethod {
  final String? id;
  final String? tenantId;
  final String name;
  final String? code;
  final bool isActive;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  const PaymentMethod({
    this.id,
    this.tenantId,
    required this.name,
    this.code,
    this.isActive = true,
    this.createdAt,
    this.updatedAt,
  });

  factory PaymentMethod.fromJson(Map<String, dynamic> json) {
    return PaymentMethod(
      id: json['id']?.toString(),
      tenantId: json['tenant_id']?.toString(),
      name: json['name']?.toString() ?? '',
      code: json['code']?.toString(),
      isActive: _toBool(json['is_active']),
      createdAt: _parseDate(json['created_at']),
      updatedAt: _parseDate(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() => {
        if (id != null) 'id': id,
        if (tenantId != null) 'tenant_id': tenantId,
        'name': name,
        if (code != null) 'code': code,
        'is_active': isActive,
        if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
        if (updatedAt != null) 'updated_at': updatedAt!.toIso8601String(),
      };

  PaymentMethod copyWith({
    String? id,
    String? tenantId,
    String? name,
    String? code,
    bool? isActive,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return PaymentMethod(
      id: id ?? this.id,
      tenantId: tenantId ?? this.tenantId,
      name: name ?? this.name,
      code: code ?? this.code,
      isActive: isActive ?? this.isActive,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }
}

bool _toBool(dynamic value) {
  if (value is bool) return value;
  if (value is int) return value != 0;
  if (value is String) return value.toLowerCase() == 'true' || value == '1';
  return false;
}

DateTime? _parseDate(dynamic value) {
  if (value == null) return null;
  if (value is DateTime) return value;
  final str = value.toString();
  if (str.isEmpty) return null;
  return DateTime.tryParse(str);
}
