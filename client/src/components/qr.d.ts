// qr.d.ts
declare module 'qr.js' {
    export interface QRCode {
        modules: boolean[][];
    }

    /**
     * Generates a QR code for a given input string.
     * 
     * @param input The string to encode in the QR code.
     * @returns An object representing the QR code with a `modules` property.
     */
    export default function qr(input: string): QRCode;
}
