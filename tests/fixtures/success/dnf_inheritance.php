<?php
interface AncA {}
interface AncB {}
interface DnfContract { public function take((AncA&AncB)|null $value): (AncA&AncB)|null; }
class AncBase implements AncA, AncB, DnfContract {
    public function take((AncB&AncA)|null $value): (AncB&AncA)|null { return $value; }
}
echo gettype((new AncBase)->take(new AncBase)), "\n";
