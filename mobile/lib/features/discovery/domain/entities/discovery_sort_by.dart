enum DiscoverySortBy {
  recommended,
  priceAsc,
  priceDesc,
  distanceAsc,
  ratingDesc,
}

extension DiscoverySortByExtension on DiscoverySortBy {
  String toParam() {
    switch (this) {
      case DiscoverySortBy.priceAsc:
        return 'PRICE_ASC';
      case DiscoverySortBy.priceDesc:
        return 'PRICE_DESC';
      case DiscoverySortBy.distanceAsc:
        return 'DISTANCE_ASC';
      case DiscoverySortBy.ratingDesc:
        return 'RATING_DESC';
      case DiscoverySortBy.recommended:
        return 'RECOMMENDED';
    }
  }
}
