<?php
class ReferenceUpdater {
    public function update(int &$value) { $value += 3; }
}
$amount = 2;
(new ReferenceUpdater)->update($amount);
echo $amount;
