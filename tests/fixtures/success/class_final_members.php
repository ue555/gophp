<?php
class FinalBase {
    final public function number() { return 9; }
}
class FinalChild extends FinalBase {}
final class FinalLeaf extends FinalChild {}
echo (new FinalLeaf)->number(), "\n";
