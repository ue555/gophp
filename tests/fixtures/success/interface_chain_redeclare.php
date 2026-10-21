<?php
function show($v) { echo json_encode($v), "\n"; }
interface ChainRoot { public function dispatch(object $event); }
interface ChainMiddle extends ChainRoot { public function dispatch(object $event, ?string $name = null): object; }
interface ChainLeaf extends ChainMiddle {}
class ChainDispatcher implements ChainLeaf {
    public function dispatch(object $event, ?string $name = null): object { return $event; }
}
$event = new stdClass();
show((new ChainDispatcher())->dispatch($event) === $event);
