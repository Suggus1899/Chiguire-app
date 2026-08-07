/// Type-safe model for the `invoices` entity.
class Invoice {
  final String? id;
  final String? tenantId;
  final String? customerId;
  final String? number;
  final String status;
  final String currency;
  final double subtotal;
  final double tax;
  final double total;
  final DateTime? issuedAt;
  final DateTime? createdAt;
  final DateTime? updatedAt;

  const Invoice({
    this.id,
    this.tenantId,
    this.customerId,
    this.number,
    this.status = 'draft',
    this.currency = 'VES',
    this.subtotal = 0,
    this.tax = 0,
    this.total = 0,
    this.issuedAt,
    this.createdAt,
    this.updatedAt,
  });

  factory Invoice.fromJson(Map<String, dynamic> json) {
    return Invoice(
      id: json['id']?.toString(),
      tenantId: json['tenant_id']?.toString(),
      customerId: json['customer_id']?.toString(),
      number: json['number']?.toString(),
      status: json['status']?.toString() ?? 'draft',
      currency: json['currency']?.toString() ?? 'VES',
      subtotal: _toDouble(json['subtotal']),
      tax: _toDouble(json['tax']),
      total: _toDouble(json['total']),
      issuedAt: _parseDate(json['issued_at']),
      createdAt: _parseDate(json['created_at']),
      updatedAt: _parseDate(json['updated_at']),
    );
  }

  Map<String, dynamic> toJson() => {
        if (id != null) 'id': id,
        if (tenantId != null) 'tenant_id': tenantId,
        if (customerId != null) 'customer_id': customerId,
        if (number != null) 'number': number,
        'status': status,
        'currency': currency,
        'subtotal': subtotal,
        'tax': tax,
        'total': total,
        if (issuedAt != null) 'issued_at': issuedAt!.toIso8601String(),
        if (createdAt != null) 'created_at': createdAt!.toIso8601String(),
        if (updatedAt != null) 'updated_at': updatedAt!.toIso8601String(),
      };

  Invoice copyWith({
    String? id,
    String? tenantId,
    String? customerId,
    String? number,
    String? status,
    String? currency,
    double? subtotal,
    double? tax,
    double? total,
    DateTime? issuedAt,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Invoice(
      id: id ?? this.id,
      tenantId: tenantId ?? this.tenantId,
      customerId: customerId ?? this.customerId,
      number: number ?? this.number,
      status: status ?? this.status,
      currency: currency ?? this.currency,
      subtotal: subtotal ?? this.subtotal,
      tax: tax ?? this.tax,
      total: total ?? this.total,
      issuedAt: issuedAt ?? this.issuedAt,
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
