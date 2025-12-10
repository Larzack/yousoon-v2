import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:cached_network_image/cached_network_image.dart';

import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/theme/app_typography.dart';

/// Écran des messages/conversations
class MessagesScreen extends ConsumerStatefulWidget {
  const MessagesScreen({super.key});

  @override
  ConsumerState<MessagesScreen> createState() => _MessagesScreenState();
}

class _MessagesScreenState extends ConsumerState<MessagesScreen> {
  // Mock data pour les conversations
  final List<Map<String, dynamic>> _conversations = [
    {
      'id': '1',
      'name': 'Le Petit Bistrot',
      'avatar': 'https://picsum.photos/seed/partner1/100',
      'lastMessage': 'Votre réservation est confirmée pour ce soir !',
      'time': '14:30',
      'unread': 2,
      'isPartner': true,
    },
    {
      'id': '2',
      'name': 'Bar L\'Éclipse',
      'avatar': 'https://picsum.photos/seed/partner2/100',
      'lastMessage': 'Merci pour votre visite, à bientôt !',
      'time': 'Hier',
      'unread': 0,
      'isPartner': true,
    },
    {
      'id': '3',
      'name': 'Support Yousoon',
      'avatar': null,
      'lastMessage': 'Comment pouvons-nous vous aider ?',
      'time': 'Mar',
      'unread': 1,
      'isPartner': false,
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        backgroundColor: AppColors.background,
        elevation: 0,
        centerTitle: true,
        title: Text('Messages', style: AppTypography.headline2),
        actions: [
          IconButton(
            icon: const Icon(Icons.edit_square, color: AppColors.textPrimary),
            onPressed: () {
              // TODO: New conversation
            },
          ),
        ],
      ),
      body: _conversations.isEmpty
          ? _buildEmptyState()
          : ListView.builder(
              itemCount: _conversations.length,
              padding: const EdgeInsets.symmetric(vertical: AppSpacing.sm),
              itemBuilder: (context, index) {
                final conversation = _conversations[index];
                return _buildConversationTile(conversation);
              },
            ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.chat_bubble_outline,
              size: 80,
              color: AppColors.inactive,
            ),
            const SizedBox(height: AppSpacing.lg),
            Text(
              'Pas de messages',
              style: AppTypography.headline2,
            ),
            const SizedBox(height: AppSpacing.sm),
            Text(
              'Vos conversations avec les partenaires\napparaîtront ici',
              style: AppTypography.bodyText.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildConversationTile(Map<String, dynamic> conversation) {
    final hasUnread = (conversation['unread'] as int) > 0;
    
    return ListTile(
      onTap: () {
        // TODO: Navigate to conversation detail
      },
      contentPadding: const EdgeInsets.symmetric(
        horizontal: AppSpacing.xl,
        vertical: AppSpacing.sm,
      ),
      leading: Stack(
        children: [
          CircleAvatar(
            radius: 28,
            backgroundColor: AppColors.surface,
            child: conversation['avatar'] != null
                ? ClipOval(
                    child: CachedNetworkImage(
                      imageUrl: conversation['avatar'],
                      width: 56,
                      height: 56,
                      fit: BoxFit.cover,
                      placeholder: (context, url) => Container(
                        color: AppColors.surface,
                      ),
                    ),
                  )
                : Icon(
                    conversation['isPartner'] 
                        ? Icons.store 
                        : Icons.support_agent,
                    color: AppColors.primary,
                  ),
          ),
          if (hasUnread)
            Positioned(
              right: 0,
              top: 0,
              child: Container(
                width: 18,
                height: 18,
                decoration: BoxDecoration(
                  color: AppColors.primary,
                  shape: BoxShape.circle,
                  border: Border.all(color: AppColors.background, width: 2),
                ),
                child: Center(
                  child: Text(
                    '${conversation['unread']}',
                    style: AppTypography.caption.copyWith(
                      color: AppColors.onPrimary,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ),
            ),
        ],
      ),
      title: Text(
        conversation['name'],
        style: AppTypography.headline3.copyWith(
          fontWeight: hasUnread ? FontWeight.bold : FontWeight.normal,
        ),
      ),
      subtitle: Padding(
        padding: const EdgeInsets.only(top: 4),
        child: Text(
          conversation['lastMessage'],
          style: AppTypography.bodyText.copyWith(
            color: hasUnread 
                ? AppColors.textPrimary 
                : AppColors.textSecondary,
            fontWeight: hasUnread ? FontWeight.w500 : FontWeight.normal,
          ),
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
      ),
      trailing: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.end,
        children: [
          Text(
            conversation['time'],
            style: AppTypography.caption.copyWith(
              color: hasUnread ? AppColors.primary : AppColors.textDisabled,
            ),
          ),
          if (hasUnread)
            Padding(
              padding: const EdgeInsets.only(top: 4),
              child: Container(
                width: 8,
                height: 8,
                decoration: const BoxDecoration(
                  color: AppColors.primary,
                  shape: BoxShape.circle,
                ),
              ),
            ),
        ],
      ),
    );
  }
}
