<?php
trait RequiresAction {
    abstract public function action();
    public function twice() { return $this->action() * 2; }
}
class ActionBase {
    use RequiresAction { twice as final; }
    public function action() { return 4; }
}
class ActionChild extends ActionBase {}
echo (new ActionChild)->twice(), "\n";
