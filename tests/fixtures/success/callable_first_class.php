<?php
function callableDouble(int $n): int { return $n * 2; }
class FirstClassBox {
    public int $count = 0;
    public function step(int $n): int { $this->count++; return $n + $this->count; }
    public static function scale(int $n): int { return $n * 3; }
}
$box = new FirstClassBox;
$function = callableDouble(...);
$method = $box->step(...);
$static = FirstClassBox::scale(...);
echo $function(3), '|', $method(4), '|', $method(4), '|', $static(3), '|', $box->count;
$runs = 0;
function countedTarget(): string { global $runs; $runs++; return 'called'; }
$lazy = countedTarget(...);
echo '|', $runs, '|', $lazy(), '|', $runs;
$creations = 0;
function makeFirstClassBox(): FirstClassBox { global $creations; $creations++; return new FirstClassBox; }
$fromReceiver = makeFirstClassBox()->step(...);
echo '|', $creations, '|', $fromReceiver(2), '|', $creations;
