/// Type-safe model for the `sellers` entity.
class Seller {
  final String? id;
  final String? tenantId;
  final String name;
  final String? email;
  final double commissionPct;
  final bool isActive;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  const Seller({
    this.id,
    this.tenantId,
    required this.name,
    this.email,
    this.commissionPct = 0,
    this.isActive = true,
    this.createdAt,
    this.updatedAt,
  });

  factory Seller.fromJson(Map<String, dynamic> json) {
    return Seller(
      id: json['id']?.toString(),
      tenantId: json['tenant_id']?.toString(),
      name: json['name']?.toString() ?? '',
      email: json['email']?.toString(),
      commissionPct: _toDouble(json['commission_pct']),
      isActive: _toBool(json['is_active']),
      createdAt: _parseDate(json['created_at']),
      updatedAt: _parseDate(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() => {
        if (id != null) 'id': id,
        if (tenantId != null) 'tenant_id': tenantId,
        'name': name,
        if (email != null) 'email': email,
        'commission_pct': commissionPct,
        'is_active': isActive,
        if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
        if (updatedAt != null) 'updated_at': updatedAt!.toIso8601String(),
      };

  Seller copyWith({
    String? id,
    String? tenantId,
    String? name,
    String? email,
    double? commissionPct,
    bool? isActive,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Seller(
      id: id ?? this.id,
      tenantId: tenantId ?? this.tenantId,
      name: name ?? this.name,
      email: email ?? this.email,
      commissionPct: commissionPct ?? this.commissionPct,
      isActive: isActive ?? this.isActive,
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
