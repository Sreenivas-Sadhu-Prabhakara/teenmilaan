import 'package:flutter/material.dart';

void main() => runApp(const TeenmilaanApp());

/// Teenmilaan — counter-side three-way match: ordered / invoiced / counted.
/// Flags any quantity or rate mismatch. Works as a 2-way match too.
class TeenmilaanApp extends StatelessWidget {
  const TeenmilaanApp({super.key});
  @override
  Widget build(BuildContext context) => MaterialApp(
        title: 'Teenmilaan',
        debugShowCheckedModeBanner: false,
        theme: ThemeData(colorSchemeSeed: const Color(0xFF6E4A8E), useMaterial3: true),
        home: const HomePage(),
      );
}

class Match {
  final String item;
  final double invoicedQty, countedQty, invoicedRate, agreedRate;
  Match(this.item, this.invoicedQty, this.countedQty, this.invoicedRate, this.agreedRate);
  double get qtyGap => invoicedQty - countedQty;
  double get rateGap => invoicedRate - agreedRate;
  bool get mismatch => qtyGap.abs() > 1e-9 || rateGap.abs() > 1e-9;
  double get valueLeak => qtyGap * agreedRate + rateGap * countedQty;
}

class HomePage extends StatefulWidget {
  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _matches = <Match>[];
  final _item = TextEditingController();
  final _iq = TextEditingController();
  final _cq = TextEditingController();
  final _ir = TextEditingController();
  final _ar = TextEditingController();

  double _n(TextEditingController c) => double.tryParse(c.text.trim()) ?? 0;

  void _add() {
    if (_item.text.trim().isEmpty) return;
    setState(() {
      _matches.insert(0, Match(_item.text.trim(), _n(_iq), _n(_cq), _n(_ir), _n(_ar)));
      _item.clear(); _iq.clear(); _cq.clear(); _ir.clear(); _ar.clear();
    });
  }

  @override
  Widget build(BuildContext context) {
    final leak = _matches.where((m) => m.mismatch).fold(0.0, (s, m) => s + m.valueLeak);
    final mis = _matches.where((m) => m.mismatch).length;
    return Scaffold(
      appBar: AppBar(
        title: const Text('Teenmilaan · 3-way match'),
        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
      ),
      body: Column(children: [
        Container(
          width: double.infinity,
          color: Theme.of(context).colorScheme.primaryContainer,
          padding: const EdgeInsets.all(16),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text('$mis mismatch${mis == 1 ? '' : 'es'} · ₹${leak.toStringAsFixed(2)} at risk',
                style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
          ]),
        ),
        Padding(padding: const EdgeInsets.all(12), child: Column(children: [
          TextField(controller: _item, decoration: const InputDecoration(labelText: 'Item', border: OutlineInputBorder())),
          const SizedBox(height: 8),
          Row(children: [
            Expanded(child: _n2(_iq, 'Invoiced qty')), const SizedBox(width: 8),
            Expanded(child: _n2(_cq, 'Counted qty')),
          ]),
          const SizedBox(height: 8),
          Row(children: [
            Expanded(child: _n2(_ir, 'Invoiced rate')), const SizedBox(width: 8),
            Expanded(child: _n2(_ar, 'Agreed rate')), const SizedBox(width: 8),
            FilledButton(onPressed: _add, child: const Text('Check')),
          ]),
        ])),
        const Divider(),
        Expanded(child: ListView.builder(
          itemCount: _matches.length,
          itemBuilder: (_, i) {
            final m = _matches[i];
            return ListTile(
              leading: Icon(m.mismatch ? Icons.warning_amber : Icons.check_circle,
                  color: m.mismatch ? Colors.red : Colors.green),
              title: Text(m.item),
              subtitle: Text('qty gap ${m.qtyGap.toStringAsFixed(0)} · rate gap ${m.rateGap.toStringAsFixed(2)}'),
              trailing: Text(m.mismatch ? '₹${m.valueLeak.toStringAsFixed(2)}' : 'ok',
                  style: const TextStyle(fontWeight: FontWeight.bold)),
            );
          },
        )),
      ]),
    );
  }

  Widget _n2(TextEditingController c, String label) => TextField(
        controller: c, keyboardType: const TextInputType.numberWithOptions(decimal: true),
        decoration: InputDecoration(labelText: label, border: const OutlineInputBorder()),
      );
}
