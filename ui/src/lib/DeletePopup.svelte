<script lang="ts">
  export var message = "Сигурни ли сте?";
  export var subMessage = "";
  export var confirmText = "Изтрий";
  export var cancelText = "Отмени";

  // let { cancelAction, confirmAction } = $props();
  export var cancelAction = () => {};
  export var confirmAction = () => {};
  export var onComplete = () => {};

  function internalCacnelAction() {
    cancelAction();
    onComplete();
  }

  function internalConfirmAction() {
    confirmAction();
    onComplete();
  }
</script>

<div
  class="fixed inset-0 z-40"
  style="background-color: rgba(107, 114, 128, 0.5);"
></div>

<div class="fixed inset-0 flex items-center justify-center z-50">
  <div class="bg-white rounded-lg p-6 shadow-lg w-96">
    <h2 class="text-lg font-semibold text-gray-800 mb-4">
      {message}
    </h2>
    <p class="text-gray-600 mb-6">
      {subMessage}
    </p>
    <div class="flex justify-end space-x-4">
      <button
        class="py-2 px-4 bg-gray-300 text-gray-700 rounded hover:bg-gray-400"
        onclick={internalCacnelAction}
      >
        {cancelText}
      </button>
      <button
        class="py-2 px-4 bg-red-500 text-white rounded hover:bg-red-600"
        onclick={internalConfirmAction}
      >
        {confirmText}
      </button>
    </div>
  </div>
</div>
