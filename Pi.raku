# Calculate pi using BBP formula with FatRat for arbitrary accuracy
say [+] (^100).map({
    FatRat.new(1, 16**$_) * (
        FatRat.new(4, 8*$_ + 1) -
        FatRat.new(2, 8*$_ + 4) -
        FatRat.new(1, 8*$_ + 5) -
        FatRat.new(1, 8*$_ + 6)
    )
});
