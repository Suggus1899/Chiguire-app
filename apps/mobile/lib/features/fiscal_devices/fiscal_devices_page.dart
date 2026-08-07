import 'package:flutter/material.dart';
import '../../core/api.dart';
import '../../widgets/app_scaffold.dart';
import '../../widgets/empty_state.dart';
import '../../widgets/loading_indicator.dart';

class FiscalDevicesPage extends StatefulWidget {
  const FiscalDevicesPage({super.key});

  @override
  State<FiscalDevicesPage> createState() => _FiscalDevicesPageState();
}

class _FiscalDevicesPageState extends State<FiscalDevicesPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  List<Map<String, dynamic>> _devices = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadDevices();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadDevices() async {
    setState(() => _isLoading = true);
    try {
      _devices = await fiscalDevicesApi.list();
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _showDeviceForm([Map<String, dynamic>? existing]) async {
    final name = TextEditingController(text: existing?['name']?.toString() ?? '');
    final model = TextEditingController(text: existing?['model']?.toString() ?? '');
    final serial = TextEditingController(text: existing?['serial']?.toString() ?? '');
    final branch = TextEditingController(text: existing?['branch_id']?.toString() ?? '');
    String type = existing?['type']?.toString() ?? 'printer';
    final ok = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setS) => Padding(
          padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom, left: 20, right: 20, top: 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(existing == null ? 'Nuevo dispositivo' : 'Editar dispositivo', style: Theme.of(ctx).textTheme.titleMedium),
              const SizedBox(height: 12),
              TextField(controller: name, decoration: const InputDecoration(labelText: 'Nombre', border: OutlineInputBorder())),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: type,
                decoration: const InputDecoration(labelText: 'Tipo', border: OutlineInputBorder()),
                items: const [
                  DropdownMenuItem(value: 'printer', child: Text('Impresora fiscal')),
                  DropdownMenuItem(value: 'pos', child: Text('POS')),
                  DropdownMenuItem(value: 'software', child: Text('Software')),
                ],
                onChanged: (v) => setS(() => type = v ?? 'printer'),
              ),
              const SizedBox(height: 12),
              TextField(controller: model, decoration: const InputDecoration(labelText: 'Modelo', border: OutlineInputBorder())),
              const SizedBox(height: 12),
              TextField(controller: serial, decoration: const InputDecoration(labelText: 'Serial', border: OutlineInputBorder())),
              const SizedBox(height: 12),
              TextField(controller: branch, decoration: const InputDecoration(labelText: 'ID Sucursal', border: OutlineInputBorder())),
              const SizedBox(height: 16),
              FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Guardar')),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
    if (ok != true) return;
    final data = {'name': name.text, 'type': type, 'model': model.text, 'serial': serial.text, 'branch_id': branch.text};
    try {
      if (existing != null) {
        await fiscalDevicesApi.update(existing['id']?.toString() ?? '', data);
      } else {
        await fiscalDevicesApi.create(data);
      }
      _loadDevices();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    return AppScaffold(
      title: 'Dispositivos Fiscales',
      currentIndex: 4,
      child: Column(
        children: [
          TabBar(
            controller: _tabController,
            tabs: const [Tab(text: 'Dispositivos'), Tab(text: 'Secuencias'), Tab(text: 'Contingencia')],
          ),
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                _devicesTab(),
                _sequencesTab(),
                _contingencyTab(),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _devicesTab() {
    if (_isLoading) return const LoadingIndicator();
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(12),
          child: Align(
            alignment: Alignment.centerRight,
            child: FloatingActionButton.extended(
              onPressed: () => _showDeviceForm(),
              icon: const Icon(Icons.add),
              label: const Text('Dispositivo'),
            ),
          ),
        ),
        Expanded(
          child: _devices.isEmpty
              ? const EmptyState(title: 'Sin dispositivos', icon: Icons.print_outlined)
              : ListView.builder(
                  itemCount: _devices.length,
                  itemBuilder: (context, i) {
                    final d = _devices[i];
                    return ListTile(
                      leading: const Icon(Icons.print),
                      title: Text(d['name']?.toString() ?? ''),
                      subtitle: Text('${d['type']?.toString() ?? ''} · ${d['serial']?.toString() ?? ''}'),
                      onTap: () => _showDeviceForm(d),
                    );
                  },
                ),
        ),
      ],
    );
  }

  Widget _sequencesTab() {
    if (_isLoading) return const LoadingIndicator();
    if (_devices.isEmpty) return const EmptyState(title: 'Primero crea un dispositivo', icon: Icons.list);
    return _SequencesTab(deviceId: _devices.first['id']?.toString() ?? '');
  }

  Widget _contingencyTab() {
    if (_isLoading) return const LoadingIndicator();
    if (_devices.isEmpty) return const EmptyState(title: 'Primero crea un dispositivo', icon: Icons.warning_amber);
    return _ContingencyTab(deviceId: _devices.first['id']?.toString() ?? '');
  }
}

class _SequencesTab extends StatefulWidget {
  final String deviceId;
  const _SequencesTab({required this.deviceId});

  @override
  State<_SequencesTab> createState() => _SequencesTabState();
}

class _SequencesTabState extends State<_SequencesTab> {
  List<Map<String, dynamic>> _sequences = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _isLoading = true);
    try {
      _sequences = await fiscalDevicesApi.listSequences(widget.deviceId);
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _showForm() async {
    String docType = 'invoice';
    final prefix = TextEditingController();
    final suffix = TextEditingController();
    final lastSeq = TextEditingController(text: '0');
    final ok = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setS) => Padding(
          padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom, left: 20, right: 20, top: 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('Nueva secuencia', style: Theme.of(ctx).textTheme.titleMedium),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: docType,
                decoration: const InputDecoration(labelText: 'Tipo doc', border: OutlineInputBorder()),
                items: const [
                  DropdownMenuItem(value: 'invoice', child: Text('Factura')),
                  DropdownMenuItem(value: 'credit_note', child: Text('Nota crédito')),
                  DropdownMenuItem(value: 'debit_note', child: Text('Nota débito')),
                ],
                onChanged: (v) => setS(() => docType = v ?? 'invoice'),
              ),
              const SizedBox(height: 12),
              TextField(controller: prefix, decoration: const InputDecoration(labelText: 'Prefijo', border: OutlineInputBorder())),
              const SizedBox(height: 12),
              TextField(controller: suffix, decoration: const InputDecoration(labelText: 'Sufijo', border: OutlineInputBorder())),
              const SizedBox(height: 12),
              TextField(controller: lastSeq, decoration: const InputDecoration(labelText: 'Último n°', border: OutlineInputBorder()), keyboardType: TextInputType.number),
              const SizedBox(height: 16),
              FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Guardar')),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
    if (ok != true) return;
    try {
      await fiscalDevicesApi.createSequence(widget.deviceId, {
        'doc_type': docType,
        'prefix': prefix.text,
        'suffix': suffix.text,
        'last_seq': int.tryParse(lastSeq.text) ?? 0,
      });
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) return const LoadingIndicator();
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(12),
          child: Align(alignment: Alignment.centerRight, child: FloatingActionButton.extended(onPressed: _showForm, icon: const Icon(Icons.add), label: const Text('Secuencia'))),
        ),
        Expanded(
          child: _sequences.isEmpty
              ? const EmptyState(title: 'Sin secuencias', icon: Icons.list)
              : ListView.builder(
                  itemCount: _sequences.length,
                  itemBuilder: (context, i) {
                    final s = _sequences[i];
                    return ListTile(
                      leading: const Icon(Icons.format_list_numbered),
                      title: Text(s['doc_type']?.toString() ?? ''),
                      subtitle: Text('Prefijo: ${s['prefix']?.toString() ?? ''} · Último: ${s['last_seq']?.toString() ?? '0'}'),
                    );
                  },
                ),
        ),
      ],
    );
  }
}

class _ContingencyTab extends StatefulWidget {
  final String deviceId;
  const _ContingencyTab({required this.deviceId});

  @override
  State<_ContingencyTab> createState() => _ContingencyTabState();
}

class _ContingencyTabState extends State<_ContingencyTab> {
  List<Map<String, dynamic>> _items = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    setState(() => _isLoading = true);
    try {
      _items = await fiscalDevicesApi.listContingency(widget.deviceId);
    } catch (_) {}
    if (mounted) setState(() => _isLoading = false);
  }

  Future<void> _showForm() async {
    String docType = 'invoice';
    final prefix = TextEditingController();
    final start = TextEditingController(text: '1');
    final end = TextEditingController(text: '100');
    final ok = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setS) => Padding(
          padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom, left: 20, right: 20, top: 20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('Nueva contingencia', style: Theme.of(ctx).textTheme.titleMedium),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                value: docType,
                decoration: const InputDecoration(labelText: 'Tipo doc', border: OutlineInputBorder()),
                items: const [
                  DropdownMenuItem(value: 'invoice', child: Text('Factura')),
                  DropdownMenuItem(value: 'credit_note', child: Text('Nota crédito')),
                  DropdownMenuItem(value: 'debit_note', child: Text('Nota débito')),
                ],
                onChanged: (v) => setS(() => docType = v ?? 'invoice'),
              ),
              const SizedBox(height: 12),
              TextField(controller: prefix, decoration: const InputDecoration(labelText: 'Prefijo', border: OutlineInputBorder())),
              const SizedBox(height: 12),
              TextField(controller: start, decoration: const InputDecoration(labelText: 'N° inicial', border: OutlineInputBorder()), keyboardType: TextInputType.number),
              const SizedBox(height: 12),
              TextField(controller: end, decoration: const InputDecoration(labelText: 'N° final', border: OutlineInputBorder()), keyboardType: TextInputType.number),
              const SizedBox(height: 16),
              FilledButton(onPressed: () => Navigator.pop(ctx, true), child: const Text('Guardar')),
              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
    if (ok != true) return;
    try {
      await fiscalDevicesApi.createContingency(widget.deviceId, {
        'doc_type': docType,
        'prefix': prefix.text,
        'start_seq': int.tryParse(start.text) ?? 1,
        'end_seq': int.tryParse(end.text) ?? 100,
      });
      _load();
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) return const LoadingIndicator();
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.all(12),
          child: Align(alignment: Alignment.centerRight, child: FloatingActionButton.extended(onPressed: _showForm, icon: const Icon(Icons.add), label: const Text('Contingencia'))),
        ),
        Expanded(
          child: _items.isEmpty
              ? const EmptyState(title: 'Sin contingencias', icon: Icons.warning_amber)
              : ListView.builder(
                  itemCount: _items.length,
                  itemBuilder: (context, i) {
                    final c = _items[i];
                    return ListTile(
                      leading: const Icon(Icons.warning_amber),
                      title: Text(c['doc_type']?.toString() ?? ''),
                      subtitle: Text('Prefijo: ${c['prefix']?.toString() ?? ''} · ${c['start_seq']}-${c['end_seq']}'),
                    );
                  },
                ),
        ),
      ],
    );
  }
}
