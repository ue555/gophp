<?php
$prices = ['りんご' => 150, 'パン' => 200, '牛乳' => 250];
$total = 0;

foreach ($prices as $name => $price) {
    echo $name, ': ', $price, "円\n";
    $total += $price;
}

echo '合計: ', $total, "円\n";

if ($total >= 500) {
    echo "500円以上のお買い物です。\n";
}
