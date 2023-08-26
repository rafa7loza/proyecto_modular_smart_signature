import React from 'react';
import { Document, Page, Image, StyleSheet } from '@react-pdf/renderer';

const styles = StyleSheet.create({
    page: {
      width: '1920px', // HD width
      height: '1080px', // HD height
      flexDirection: 'row',
      backgroundColor: '#FFFFFF',
      alignItems: 'center',
      justifyContent: 'center',
    },
    image: {
      width: '100%',
      height: '100%',
    },
  });

const PDFViewer = ({ base64Img, isOpen, onClose }) => {
    return (
        <div>
            {isOpen && (
                <div className="modal">
                    <Document>
                        <Page size={{ width: '800px', height: '800px' }} style={styles.page}>
                            <Image src={`data:image/jpeg;base64,${base64Img}`} style={styles.image} />
                        </Page>
                    </Document>
                    <button onClick={onClose}>Close</button>
                </div>
            )}
        </div>
    );
};

export default PDFViewer;
