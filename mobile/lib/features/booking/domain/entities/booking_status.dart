enum BookingStatus {
  draft,
  pendingPayment,
  paymentEvidenceUploaded,
  paymentDeclined,
  disputeOpened,
  disputeResolvedRenter,
  disputeResolvedOwner,
  confirmed,
  awaitingFleetReadiness,
  pendingAdditionalPayment,
  readyForPickup,
  active,
  returnInitiated,
  pendingFinalPayment,
  completed,
  rated,
  cancelled;

  static BookingStatus fromString(String val) {
    switch (val.toUpperCase()) {
      case 'DRAFT':
        return BookingStatus.draft;
      case 'PENDING_PAYMENT':
        return BookingStatus.pendingPayment;
      case 'PAYMENT_EVIDENCE_UPLOADED':
        return BookingStatus.paymentEvidenceUploaded;
      case 'PAYMENT_DECLINED':
        return BookingStatus.paymentDeclined;
      case 'DISPUTE_OPENED':
        return BookingStatus.disputeOpened;
      case 'DISPUTE_RESOLVED_RENTER':
        return BookingStatus.disputeResolvedRenter;
      case 'DISPUTE_RESOLVED_OWNER':
        return BookingStatus.disputeResolvedOwner;
      case 'CONFIRMED':
        return BookingStatus.confirmed;
      case 'AWAITING_FLEET_READINESS':
        return BookingStatus.awaitingFleetReadiness;
      case 'PENDING_ADDITIONAL_PAYMENT':
        return BookingStatus.pendingAdditionalPayment;
      case 'READY_FOR_PICKUP':
        return BookingStatus.readyForPickup;
      case 'ACTIVE':
        return BookingStatus.active;
      case 'RETURN_INITIATED':
        return BookingStatus.returnInitiated;
      case 'PENDING_FINAL_PAYMENT':
        return BookingStatus.pendingFinalPayment;
      case 'COMPLETED':
        return BookingStatus.completed;
      case 'RATED':
        return BookingStatus.rated;
      case 'CANCELLED':
        return BookingStatus.cancelled;
      default:
        return BookingStatus.pendingPayment;
    }
  }

  String toApiString() {
    switch (this) {
      case BookingStatus.draft:
        return 'DRAFT';
      case BookingStatus.pendingPayment:
        return 'PENDING_PAYMENT';
      case BookingStatus.paymentEvidenceUploaded:
        return 'PAYMENT_EVIDENCE_UPLOADED';
      case BookingStatus.paymentDeclined:
        return 'PAYMENT_DECLINED';
      case BookingStatus.disputeOpened:
        return 'DISPUTE_OPENED';
      case BookingStatus.disputeResolvedRenter:
        return 'DISPUTE_RESOLVED_RENTER';
      case BookingStatus.disputeResolvedOwner:
        return 'DISPUTE_RESOLVED_OWNER';
      case BookingStatus.confirmed:
        return 'CONFIRMED';
      case BookingStatus.awaitingFleetReadiness:
        return 'AWAITING_FLEET_READINESS';
      case BookingStatus.pendingAdditionalPayment:
        return 'PENDING_ADDITIONAL_PAYMENT';
      case BookingStatus.readyForPickup:
        return 'READY_FOR_PICKUP';
      case BookingStatus.active:
        return 'ACTIVE';
      case BookingStatus.returnInitiated:
        return 'RETURN_INITIATED';
      case BookingStatus.pendingFinalPayment:
        return 'PENDING_FINAL_PAYMENT';
      case BookingStatus.completed:
        return 'COMPLETED';
      case BookingStatus.rated:
        return 'RATED';
      case BookingStatus.cancelled:
        return 'CANCELLED';
    }
  }
}
