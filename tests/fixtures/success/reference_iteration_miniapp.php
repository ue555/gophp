<?php
class ReferenceLedger implements Iterator, Countable {
    private array $rows = [2, 3];
    private int $position = 0;
    public function count(): int { return count($this->rows); }
    public function rewind(): void { $this->position = 0; }
    public function valid(): bool { return isset($this->rows[$this->position]); }
    public function key(): int { return $this->position; }
    public function current(): int { return $this->rows[$this->position]; }
    public function next(): void { $this->position++; }
}
function addTotal(int &$total, int $amount): void { $total += $amount; }
$ledger = new ReferenceLedger;
$total = 0;
foreach ($ledger as $amount) addTotal($total, $amount);
echo count($ledger), '|', $total;
