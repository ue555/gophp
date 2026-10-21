<?php
function loadedVoid(): void { return; }
function loadedNever(): never { throw new Exception('loaded'); }
