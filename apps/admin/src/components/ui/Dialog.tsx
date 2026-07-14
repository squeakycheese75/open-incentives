import { useEffect, useRef, type ReactNode } from "react";

interface DialogProps {
  open: boolean;
  onClose: () => void;
  title: string;
  children: ReactNode;
}

export function Dialog({ open, onClose, title, children }: DialogProps) {
  const ref = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    if (open && !el.open) {
      el.showModal();
    } else if (!open && el.open) {
      el.close();
    }
  }, [open]);

  return (
    <dialog
      ref={ref}
      onClose={onClose}
      onCancel={onClose}
      className="w-full max-w-md rounded-lg border border-gray-200 p-0 shadow-lg backdrop:bg-black/30"
    >
      <div className="flex items-center justify-between border-b border-gray-200 px-5 py-4">
        <h2 className="text-base font-semibold text-gray-900">{title}</h2>
        <button onClick={onClose} className="text-gray-400 hover:text-gray-600" aria-label="Close">
          &times;
        </button>
      </div>
      <div className="p-5">{children}</div>
    </dialog>
  );
}
