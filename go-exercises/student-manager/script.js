function openModal(modalId) {
            document.getElementById(modalId).classList.remove('hidden');
        }

        function closeModal(modalId) {
            document.getElementById(modalId).classList.add('hidden');
        }

        function simulateAction(actionName) {
            const consoleLog = document.getElementById('consoleLogText');
            consoleLog.textContent = `Executed: ${actionName} [UI simulation mockup mode]`;
            consoleLog.classList.remove('text-emerald-400');
            consoleLog.classList.add('text-amber-400');
            setTimeout(() => {
                consoleLog.classList.remove('text-amber-400');
                consoleLog.classList.add('text-emerald-400');
            }, 1000);
        }

        function handleModalSubmit(event, successMessage) {
            event.preventDefault();
            const consoleLog = document.getElementById('consoleLogText');
            consoleLog.textContent = successMessage;
            
            // Close all modals
            document.querySelectorAll('[id$="Modal"]').forEach(el => el.classList.add('hidden'));
            
            // Show alert notice
            alert(successMessage + "\n(Note: This is a non-functional frontend UI mockup connected to your Go backend logic)");
        }
