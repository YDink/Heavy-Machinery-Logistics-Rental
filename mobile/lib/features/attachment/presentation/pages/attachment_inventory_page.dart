import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:dink_mobile/l10n/app_localizations.dart';
import 'package:dink_mobile/features/attachment/domain/entities/attachment_inventory_item.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_bloc.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_event.dart';
import 'package:dink_mobile/features/attachment/presentation/bloc/attachment_state.dart';
import 'package:dink_mobile/features/attachment/presentation/pages/attachment_registration_page.dart';

class AttachmentInventoryPage extends StatefulWidget {
  const AttachmentInventoryPage({super.key});

  @override
  State<AttachmentInventoryPage> createState() => _AttachmentInventoryPageState();
}

class _AttachmentInventoryPageState extends State<AttachmentInventoryPage> {
  @override
  void initState() {
    super.initState();
    context.read<AttachmentBloc>().add(const LoadAttachmentInventoryEvent());
  }

  @override
  Widget build(BuildContext context) {
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.attachmentInventoryTitle),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            tooltip: l10n.registerAttachmentTitle,
            onPressed: () {
              Navigator.of(context).push(
                MaterialPageRoute(builder: (_) => const AttachmentRegistrationPage()),
              );
            },
          ),
        ],
      ),
      body: BlocConsumer<AttachmentBloc, AttachmentState>(
        listener: (context, state) {
          if (state.status == AttachmentStatusState.linked) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(l10n.attachmentLinkedSuccess),
                backgroundColor: Colors.green,
              ),
            );
          } else if (state.status == AttachmentStatusState.error) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.errorMessage ?? l10n.errorGeneric),
                backgroundColor: Colors.red,
              ),
            );
          }
        },
        builder: (context, state) {
          if (state.status == AttachmentStatusState.loading && state.inventory.isEmpty) {
            return const Center(child: CircularProgressIndicator.adaptive());
          }

          if (state.inventory.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.handyman_outlined, size: 64, color: Colors.blueGrey),
                  const SizedBox(height: 16),
                  Text(
                    l10n.emptyAttachmentInventory,
                    style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 12),
                  ElevatedButton.icon(
                    icon: const Icon(Icons.add),
                    label: Text(l10n.registerAttachmentTitle),
                    onPressed: () {
                      Navigator.of(context).push(
                        MaterialPageRoute(builder: (_) => const AttachmentRegistrationPage()),
                      );
                    },
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () async {
              context.read<AttachmentBloc>().add(const LoadAttachmentInventoryEvent());
            },
            child: ListView.separated(
              padding: const EdgeInsets.all(16),
              itemCount: state.inventory.length,
              separatorBuilder: (context, index) => const SizedBox(height: 12),
              itemBuilder: (context, index) {
                final item = state.inventory[index];
                return _buildAttachmentCard(context, item, l10n);
              },
            ),
          );
        },
      ),
    );
  }

  Widget _buildAttachmentCard(BuildContext context, AttachmentInventoryItem item, AppLocalizations l10n) {
    final dailyRateFormatted = (item.standaloneDailyRateMinor / 100).toStringAsFixed(2);
    final isActive = item.status == 'ACTIVE';

    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    '${item.make} ${item.model}',
                    style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16),
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: isActive ? Colors.green.shade50 : Colors.orange.shade50,
                    borderRadius: BorderRadius.circular(6),
                    border: Border.all(color: isActive ? Colors.green : Colors.orange),
                  ),
                  child: Text(
                    isActive ? l10n.attachmentStatusActive : item.status,
                    style: TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.bold,
                      color: isActive ? Colors.green.shade800 : Colors.orange.shade800,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 6),
            Text(
              '${item.attachmentType} • ${item.subTypeSize}',
              style: TextStyle(fontSize: 13, color: Colors.grey.shade700),
            ),
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  '$dailyRateFormatted ETB / ${l10n.perDay}',
                  style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15, color: Colors.blueGrey),
                ),
                if (item.isLinkedToFleet)
                  Chip(
                    label: Text(l10n.bundledWithHost, style: const TextStyle(fontSize: 11)),
                    backgroundColor: const Color(0xFFE3F2FD),
                  ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
